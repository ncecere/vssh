package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vssh/pkg/types"

	"github.com/sirupsen/logrus"
)

// Client handles SSH client operations
type Client struct {
	config *types.Config
	logger *logrus.Logger
}

// NewClient creates a new SSH client
func NewClient(config *types.Config, logger *logrus.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

// SSHOptions represents SSH command options
type SSHOptions struct {
	Port            string
	IdentityFile    string
	CertificateFile string
	IPv4            bool
	IPv6            bool
	Verbose         bool
	Debug           bool
	ExtraArgs       []string
}

// Connect executes SSH connection with the signed certificate
func (c *Client) Connect(target *SSHTarget, certPath string, options *SSHOptions, command []string) error {
	// Build SSH command arguments
	args := []string{}

	// Add port if specified
	if options.Port != "" {
		args = append(args, "-p", options.Port)
	}

	// Add certificate file
	if certPath != "" {
		args = append(args, "-o", fmt.Sprintf("CertificateFile=%s", certPath))
	}

	// Add identity file if specified
	if options.IdentityFile != "" {
		args = append(args, "-i", options.IdentityFile)
	}

	// Add IP version flags
	if options.IPv4 {
		args = append(args, "-4")
	}
	if options.IPv6 {
		args = append(args, "-6")
	}

	// Add verbose/debug flags
	if options.Verbose {
		args = append(args, "-v")
	}
	if options.Debug {
		args = append(args, "-vvv")
	}

	// Add extra SSH options for certificate-based authentication
	args = append(args, "-o", "PreferredAuthentications=publickey")
	args = append(args, "-o", "PubkeyAuthentication=yes")

	// Add any extra arguments
	args = append(args, options.ExtraArgs...)

	// Add the target (user@hostname)
	sshTarget := fmt.Sprintf("%s@%s", target.Username, target.Hostname)
	args = append(args, sshTarget)

	// Add command if specified
	if len(command) > 0 {
		args = append(args, command...)
	}

	c.logger.Debugf("Executing SSH command: ssh %s", strings.Join(args, " "))

	// Execute SSH command
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables if needed
	cmd.Env = os.Environ()

	// Execute the command
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// SSH command failed, return the exit code
			return fmt.Errorf("SSH connection failed with exit code %d", exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute SSH command: %w", err)
	}

	return nil
}

// ParseSSHArgs parses SSH command line arguments and extracts options
func ParseSSHArgs(args []string) (*SSHOptions, []string, error) {
	options := &SSHOptions{}
	var command []string

	// Simple implementation - just find the target and treat everything else as command
	if len(args) == 0 {
		return options, command, fmt.Errorf("no arguments provided")
	}

	// For now, assume first argument is always the target
	// and everything after is the command
	if len(args) > 1 {
		command = args[1:]
	}

	return options, command, nil
}

// GetPrivateKeyPath returns the private key path for the certificate
func (c *Client) GetPrivateKeyPath(username string) (string, error) {
	// Check if user has specific configuration
	if userConfig, exists := c.config.Users[username]; exists {
		return userConfig.PrivateKey, nil
	}

	// Use default key path
	keyPath := filepath.Join(c.config.SSH.KeyDirectory, "id_rsa")

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

// ValidateSSHBinary checks if SSH binary is available
func (c *Client) ValidateSSHBinary() error {
	_, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("SSH binary not found in PATH. Please install OpenSSH client")
	}
	return nil
}
