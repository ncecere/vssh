package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"vssh/internal/vault"
	"vssh/pkg/types"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

// Authenticator handles Vault authentication
type Authenticator struct {
	client *vault.Client
	config *types.VaultConfig
	logger *logrus.Logger
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(client *vault.Client, config *types.VaultConfig, logger *logrus.Logger) *Authenticator {
	return &Authenticator{
		client: client,
		config: config,
		logger: logger,
	}
}

// EnsureAuthenticated ensures the client has a valid token, prompting for authentication if needed
func (a *Authenticator) EnsureAuthenticated() error {
	// First, try to load existing token
	if err := a.client.LoadTokenFromFile(); err != nil {
		a.logger.Debugf("Could not load token from file: %v", err)
	}

	// Check if current token is valid
	if a.client.IsTokenValid() {
		a.logger.Debug("Using existing valid token")
		return nil
	}

	a.logger.Info("No valid token found, authentication required")

	// Determine authentication method
	authMethod := types.AuthMethod(a.config.AuthMethod)

	// If no auth method configured, prompt user to choose
	if authMethod == "" || !authMethod.IsValid() {
		var err error
		authMethod, err = a.promptForAuthMethod()
		if err != nil {
			return fmt.Errorf("failed to get authentication method: %w", err)
		}
	}

	// Authenticate using the selected method
	if err := a.authenticate(authMethod); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Save the new token
	if err := a.client.SaveTokenToFile(); err != nil {
		a.logger.Warnf("Failed to save token to file: %v", err)
		// Don't fail here, token is still valid in memory
	}

	a.logger.Info("Authentication successful")
	return nil
}

// promptForAuthMethod prompts the user to choose an authentication method
func (a *Authenticator) promptForAuthMethod() (types.AuthMethod, error) {
	fmt.Println("Please choose an authentication method:")
	fmt.Println("1. Token")
	fmt.Println("2. Username/Password")
	fmt.Println("3. LDAP")
	fmt.Println("4. OIDC")
	fmt.Print("Enter your choice (1-4): ")

	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		return types.AuthMethodToken, nil
	case "2":
		return types.AuthMethodUserPass, nil
	case "3":
		return types.AuthMethodLDAP, nil
	case "4":
		return types.AuthMethodOIDC, nil
	default:
		return "", fmt.Errorf("invalid choice: %s", choice)
	}
}

// authenticate performs authentication using the specified method
func (a *Authenticator) authenticate(method types.AuthMethod) error {
	switch method {
	case types.AuthMethodToken:
		return a.authenticateToken()
	case types.AuthMethodUserPass:
		return a.authenticateUserPass()
	case types.AuthMethodLDAP:
		return a.authenticateLDAP()
	case types.AuthMethodOIDC:
		return a.authenticateOIDC()
	default:
		return fmt.Errorf("unsupported authentication method: %s", method)
	}
}

// authenticateToken prompts for a token and sets it
func (a *Authenticator) authenticateToken() error {
	fmt.Print("Enter Vault token: ")

	// Read token securely (hidden input)
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading token: %w", err)
	}
	fmt.Println() // Add newline after hidden input

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Set token and validate
	a.client.SetToken(token)
	if !a.client.IsTokenValid() {
		return fmt.Errorf("invalid token provided")
	}

	return nil
}

// authenticateUserPass performs username/password authentication
func (a *Authenticator) authenticateUserPass() error {
	reader := bufio.NewReader(os.Stdin)

	// Get username
	username := a.config.UserPass.Username
	if username == "" {
		fmt.Print("Username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading username: %w", err)
		}
		username = strings.TrimSpace(input)
	}

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Get password
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading password: %w", err)
	}
	fmt.Println() // Add newline after hidden input

	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Perform authentication
	mount := a.config.UserPass.Mount
	if mount == "" {
		mount = "userpass"
	}

	path := fmt.Sprintf("auth/%s/login/%s", mount, username)
	data := map[string]interface{}{
		"password": password,
	}

	secret, err := a.client.GetClient().Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("userpass authentication failed: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("no authentication data returned")
	}

	// Set the token
	a.client.SetToken(secret.Auth.ClientToken)
	return nil
}

// authenticateLDAP performs LDAP authentication
func (a *Authenticator) authenticateLDAP() error {
	reader := bufio.NewReader(os.Stdin)

	// Get username
	username := a.config.LDAP.Username
	if username == "" {
		fmt.Print("LDAP Username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading username: %w", err)
		}
		username = strings.TrimSpace(input)
	}

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Get password
	fmt.Print("LDAP Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading password: %w", err)
	}
	fmt.Println() // Add newline after hidden input

	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Perform authentication
	mount := a.config.LDAP.Mount
	if mount == "" {
		mount = "ldap"
	}

	path := fmt.Sprintf("auth/%s/login/%s", mount, username)
	data := map[string]interface{}{
		"password": password,
	}

	secret, err := a.client.GetClient().Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("LDAP authentication failed: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("no authentication data returned")
	}

	// Set the token
	a.client.SetToken(secret.Auth.ClientToken)
	return nil
}

// authenticateOIDC performs OIDC authentication
func (a *Authenticator) authenticateOIDC() error {
	mount := a.config.OIDC.Mount
	if mount == "" {
		mount = "oidc"
	}

	role := a.config.OIDC.Role
	if role == "" {
		return fmt.Errorf("OIDC role not configured")
	}

	fmt.Printf("Starting OIDC authentication for role: %s\n", role)
	fmt.Println("This will open a browser window for authentication...")

	// Start OIDC auth
	path := fmt.Sprintf("auth/%s/oidc/auth_url", mount)
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": "http://localhost:8250/oidc/callback",
	}

	secret, err := a.client.GetClient().Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("failed to get OIDC auth URL: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return fmt.Errorf("no OIDC auth URL returned")
	}

	authURL, ok := secret.Data["auth_url"].(string)
	if !ok {
		return fmt.Errorf("invalid auth URL returned")
	}

	fmt.Printf("Please visit this URL to authenticate: %s\n", authURL)
	fmt.Print("Enter the authorization code: ")

	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading authorization code: %w", err)
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return fmt.Errorf("authorization code cannot be empty")
	}

	// Complete OIDC authentication
	completePath := fmt.Sprintf("auth/%s/oidc/callback", mount)
	completeData := map[string]interface{}{
		"code":  code,
		"state": secret.Data["state"],
	}

	authSecret, err := a.client.GetClient().Logical().Write(completePath, completeData)
	if err != nil {
		return fmt.Errorf("OIDC authentication failed: %w", err)
	}

	if authSecret == nil || authSecret.Auth == nil {
		return fmt.Errorf("no authentication data returned")
	}

	// Set the token
	a.client.SetToken(authSecret.Auth.ClientToken)
	return nil
}
