package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vssh/internal/vault"
	"vssh/pkg/types"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// Signer handles SSH key signing operations
type Signer struct {
	vaultClient *vault.Client
	config      *types.Config
	logger      *logrus.Logger
}

// NewSigner creates a new SSH signer
func NewSigner(vaultClient *vault.Client, config *types.Config, logger *logrus.Logger) *Signer {
	return &Signer{
		vaultClient: vaultClient,
		config:      config,
		logger:      logger,
	}
}

// SSHTarget represents a parsed SSH connection target
type SSHTarget struct {
	Username string
	Hostname string
	Port     string
}

// ParseSSHTarget parses an SSH target string like "user@hostname" or "hostname"
func ParseSSHTarget(target string) (*SSHTarget, error) {
	sshTarget := &SSHTarget{}

	// Split on @ to separate user and host
	parts := strings.Split(target, "@")
	if len(parts) == 2 {
		sshTarget.Username = parts[0]
		sshTarget.Hostname = parts[1]
	} else if len(parts) == 1 {
		// No username specified, use current user
		currentUser := os.Getenv("USER")
		if currentUser == "" {
			return nil, fmt.Errorf("no username specified and USER environment variable not set")
		}
		sshTarget.Username = currentUser
		sshTarget.Hostname = parts[0]
	} else {
		return nil, fmt.Errorf("invalid SSH target format: %s", target)
	}

	if sshTarget.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if sshTarget.Hostname == "" {
		return nil, fmt.Errorf("hostname cannot be empty")
	}

	return sshTarget, nil
}

// GetPrivateKeyPath returns the private key path for a user
func (s *Signer) GetPrivateKeyPath(username string) (string, error) {
	// Check if user has specific configuration
	if userConfig, exists := s.config.Users[username]; exists {
		return userConfig.PrivateKey, nil
	}

	// Use default key path
	keyPath := filepath.Join(s.config.SSH.KeyDirectory, "id_rsa")

	// Expand tilde if present
	if strings.HasPrefix(keyPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %w", err)
		}
		keyPath = filepath.Join(home, keyPath[1:])
	}

	return keyPath, nil
}

// GetCertificatePath returns the path where the signed certificate should be stored
func (s *Signer) GetCertificatePath(username string) string {
	certName := fmt.Sprintf("vault_signed_%s.pub", username)
	return filepath.Join(s.config.SSH.KeyDirectory, certName)
}

// IsCertificateValid checks if an existing certificate is still valid
func (s *Signer) IsCertificateValid(certPath string) bool {
	// Check if certificate file exists
	certData, err := os.ReadFile(certPath)
	if err != nil {
		s.logger.Debugf("Certificate file not found: %s", certPath)
		return false
	}

	// Parse the certificate
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(certData)
	if err != nil {
		s.logger.Debugf("Failed to parse certificate: %v", err)
		return false
	}

	// Check if it's actually a certificate
	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		s.logger.Debug("Public key is not a certificate")
		return false
	}

	// Check if certificate is still valid (not expired)
	now := uint64(time.Now().Unix())
	if cert.ValidBefore != 0 && now >= cert.ValidBefore {
		s.logger.Debugf("Certificate expired at %d, current time %d", cert.ValidBefore, now)
		return false
	}

	// Check if certificate is not yet valid
	if cert.ValidAfter != 0 && now < cert.ValidAfter {
		s.logger.Debugf("Certificate not yet valid until %d, current time %d", cert.ValidAfter, now)
		return false
	}

	// Consider certificate valid if it has more than 5 minutes remaining
	if cert.ValidBefore != 0 {
		remaining := time.Duration(cert.ValidBefore-now) * time.Second
		if remaining < 5*time.Minute {
			s.logger.Debugf("Certificate expires soon: %v remaining", remaining)
			return false
		}
		s.logger.Debugf("Certificate is valid with %v remaining", remaining)
	}

	return true
}

// SignSSHKey signs an SSH public key using Vault
func (s *Signer) SignSSHKey(username string, publicKeyPath string) (string, error) {
	// Read the public key
	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key %s: %w", publicKeyPath, err)
	}

	// Get the vault role for this user
	// Default to using the username as the role (matches Vault CLI pattern)
	vaultRole := username

	// Allow override from user configuration
	if userConfig, exists := s.config.Users[username]; exists && userConfig.VaultRole != "" {
		vaultRole = userConfig.VaultRole
	} else if s.config.Vault.Role != "" {
		// Fallback to global role if configured (for backward compatibility)
		vaultRole = s.config.Vault.Role
	}

	s.logger.Debugf("Signing SSH key for user %s with role %s", username, vaultRole)

	// Prepare signing request
	path := fmt.Sprintf("%s/sign/%s", s.config.SSH.SigningEngine, vaultRole)
	data := map[string]interface{}{
		"public_key": string(pubKeyData),
		"ttl":        s.config.SSH.CertificateTTL.String(),
	}

	// Make the signing request to Vault
	secret, err := s.vaultClient.GetClient().Logical().Write(path, data)
	if err != nil {
		return "", fmt.Errorf("failed to sign SSH key: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no data returned from Vault SSH signing")
	}

	// Extract the signed certificate
	signedKey, ok := secret.Data["signed_key"].(string)
	if !ok {
		return "", fmt.Errorf("signed_key not found in Vault response")
	}

	s.logger.Debugf("Successfully signed SSH key for user %s", username)
	return signedKey, nil
}

// EnsureSSHCertificate ensures a valid SSH certificate exists for the user
func (s *Signer) EnsureSSHCertificate(username string) (string, error) {
	certPath := s.GetCertificatePath(username)

	// Check if we already have a valid certificate
	if s.IsCertificateValid(certPath) {
		s.logger.Debugf("Using existing valid certificate: %s", certPath)
		return certPath, nil
	}

	s.logger.Infof("Generating new SSH certificate for user: %s", username)

	// Get the private key path
	privateKeyPath, err := s.GetPrivateKeyPath(username)
	if err != nil {
		return "", fmt.Errorf("failed to get private key path: %w", err)
	}

	// Generate public key path from private key path
	publicKeyPath := privateKeyPath + ".pub"

	// Check if private key exists
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return "", fmt.Errorf("private key not found: %s. Please generate an SSH key pair first", privateKeyPath)
	}

	// Check if public key exists
	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return "", fmt.Errorf("public key not found: %s. Please generate an SSH key pair first", publicKeyPath)
	}

	// Sign the SSH key
	signedCert, err := s.SignSSHKey(username, publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to sign SSH key: %w", err)
	}

	// Ensure the SSH directory exists
	sshDir := filepath.Dir(certPath)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create SSH directory: %w", err)
	}

	// Write the signed certificate to file
	if err := os.WriteFile(certPath, []byte(signedCert), 0644); err != nil {
		return "", fmt.Errorf("failed to write certificate file: %w", err)
	}

	s.logger.Infof("SSH certificate saved to: %s", certPath)
	return certPath, nil
}
