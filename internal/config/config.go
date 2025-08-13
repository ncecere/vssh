package config

import (
	"fmt"
	"os"
	"path/filepath"

	"vssh/pkg/types"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from file and environment variables
func LoadConfig() (*types.Config, error) {
	config := &types.Config{}

	// Set defaults
	setDefaults()

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults
	}

	// Unmarshal into our config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Get home directory for default paths
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	// Vault defaults
	viper.SetDefault("vault.address", "https://vault.example.com")
	// viper.SetDefault("vault.role", "ssh-client-role")  # Removed - will use username as role
	viper.SetDefault("vault.auth_method", "token")
	viper.SetDefault("vault.token.token_path", filepath.Join(home, ".vault-token"))
	viper.SetDefault("vault.userpass.mount", "userpass")
	viper.SetDefault("vault.ldap.mount", "ldap")
	viper.SetDefault("vault.oidc.mount", "oidc")

	// SSH defaults
	viper.SetDefault("ssh.key_directory", filepath.Join(home, ".ssh"))
	viper.SetDefault("ssh.certificate_ttl", "4h")
	viper.SetDefault("ssh.signing_engine", "ssh-client-signer")

	// Debug default
	viper.SetDefault("debug", false)
}

// validateConfig validates the loaded configuration
func validateConfig(config *types.Config) error {
	// Validate Vault configuration
	if config.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}

	// vault.role is now optional - will use username as role by default

	// Validate auth method
	authMethod := types.AuthMethod(config.Vault.AuthMethod)
	if !authMethod.IsValid() {
		return fmt.Errorf("invalid auth method: %s. Supported methods: token, userpass, ldap, oidc", config.Vault.AuthMethod)
	}

	// Validate auth method specific configuration
	switch authMethod {
	case types.AuthMethodUserPass:
		if config.Vault.UserPass.Username == "" {
			// Username can be prompted at runtime, so this is not required
		}
	case types.AuthMethodLDAP:
		if config.Vault.LDAP.Username == "" {
			// Username can be prompted at runtime, so this is not required
		}
	case types.AuthMethodOIDC:
		if config.Vault.OIDC.Role == "" {
			return fmt.Errorf("vault.oidc.role is required when using oidc auth")
		}
	}

	// Validate SSH configuration
	if config.SSH.KeyDirectory == "" {
		return fmt.Errorf("ssh.key_directory is required")
	}

	// Validate certificate TTL
	if config.SSH.CertificateTTL <= 0 {
		return fmt.Errorf("ssh.certificate_ttl must be greater than 0")
	}

	// Validate user configurations
	for username, userConfig := range config.Users {
		if userConfig.PrivateKey == "" {
			return fmt.Errorf("private_key is required for user %s", username)
		}

		// Expand tilde in private key path
		if userConfig.PrivateKey[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("error getting home directory for user %s: %w", username, err)
			}
			userConfig.PrivateKey = filepath.Join(home, userConfig.PrivateKey[1:])
			config.Users[username] = userConfig
		}
	}

	return nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(configPath string) error {
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Get home directory for default paths
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	// Create default configuration content
	defaultConfig := fmt.Sprintf(`# vssh configuration file
# See https://github.com/ncecere/vssh for documentation

vault:
  address: "https://vault.example.com"
  role: "ssh-client-role"
  auth_method: "token"  # Options: token, userpass, ldap, oidc
  
  # Token authentication (default)
  token:
    token_path: "%s/.vault-token"
  
  # Username/Password authentication
  # userpass:
  #   username: "your-username"
  #   mount: "userpass"
  
  # LDAP authentication
  # ldap:
  #   username: "your-username"
  #   mount: "ldap"
  
  # OIDC authentication
  # oidc:
  #   role: "your-oidc-role"
  #   mount: "oidc"

ssh:
  key_directory: "%s/.ssh"
  certificate_ttl: "4h"
  signing_engine: "ssh-client-signer"

# Per-user SSH key configuration
users:
  # Example user configuration
  # user1:
  #   private_key: "%s/.ssh/user1_rsa"
  #   vault_role: "user1-role"  # Optional: override default vault role
  
  # user2:
  #   private_key: "%s/.ssh/user2_rsa"
  #   vault_role: "user2-role"

# Enable debug logging
debug: false
`, home, home, home, home)

	// Write the configuration file
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the configuration file path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "vssh", "config.yaml")
}
