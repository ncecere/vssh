package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vssh/internal/config"
	"vssh/pkg/types"

	"github.com/spf13/viper"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()

	// Load config without any config file (should use defaults)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test default values
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("Expected default vault address, got %s", cfg.Vault.Address)
	}

	if cfg.Vault.Role != "ssh-client-role" {
		t.Errorf("Expected default vault role, got %s", cfg.Vault.Role)
	}

	if cfg.Vault.AuthMethod != "token" {
		t.Errorf("Expected default auth method 'token', got %s", cfg.Vault.AuthMethod)
	}

	if cfg.SSH.CertificateTTL != 4*time.Hour {
		t.Errorf("Expected default certificate TTL 4h, got %v", cfg.SSH.CertificateTTL)
	}
}

func TestLoadConfig_WithCustomConfig(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	configContent := `
vault:
  address: "https://custom-vault.example.com"
  role: "custom-role"
  auth_method: "userpass"
  userpass:
    username: "testuser"
    mount: "userpass"

ssh:
  key_directory: "/tmp/ssh"
  certificate_ttl: "2h"
  ssh_engine: "ssh"

users:
  testuser:
    private_key: "/tmp/ssh/testuser_rsa"
    vault_role: "testuser-role"

debug: true
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Reset viper and set config file
	viper.Reset()
	viper.SetConfigFile(configFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test custom values
	if cfg.Vault.Address != "https://custom-vault.example.com" {
		t.Errorf("Expected custom vault address, got %s", cfg.Vault.Address)
	}

	if cfg.Vault.Role != "custom-role" {
		t.Errorf("Expected custom vault role, got %s", cfg.Vault.Role)
	}

	if cfg.Vault.AuthMethod != "userpass" {
		t.Errorf("Expected auth method 'userpass', got %s", cfg.Vault.AuthMethod)
	}

	if cfg.Vault.UserPass.Username != "testuser" {
		t.Errorf("Expected userpass username 'testuser', got %s", cfg.Vault.UserPass.Username)
	}

	if cfg.SSH.CertificateTTL != 2*time.Hour {
		t.Errorf("Expected certificate TTL 2h, got %v", cfg.SSH.CertificateTTL)
	}

	if cfg.Debug != true {
		t.Errorf("Expected debug true, got %v", cfg.Debug)
	}

	// Test user configuration
	userConfig, exists := cfg.Users["testuser"]
	if !exists {
		t.Errorf("Expected user 'testuser' to exist in config")
	}

	if userConfig.PrivateKey != "/tmp/ssh/testuser_rsa" {
		t.Errorf("Expected private key path, got %s", userConfig.PrivateKey)
	}

	if userConfig.VaultRole != "testuser-role" {
		t.Errorf("Expected vault role 'testuser-role', got %s", userConfig.VaultRole)
	}
}

func TestLoadConfig_ValidationFailure(t *testing.T) {
	// Create temporary config file with invalid configuration
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	configContent := `
vault:
  address: ""  # Invalid: empty address
  role: "test-role"
  auth_method: "token"

ssh:
  key_directory: "/tmp/ssh"
  certificate_ttl: "4h"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Reset viper and set config file
	viper.Reset()
	viper.SetConfigFile(configFile)

	_, err = config.LoadConfig()
	if err == nil {
		t.Errorf("Expected validation error for empty vault address, got nil")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := config.CreateDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error creating default config, got %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Expected config file to be created at %s", configPath)
	}

	// Check if file contains expected content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	contentStr := string(content)
	expectedStrings := []string{
		"vault:",
		"address:",
		"auth_method:",
		"ssh:",
		"certificate_ttl:",
		"users:",
	}

	for _, expected := range expectedStrings {
		if !contains(contentStr, expected) {
			t.Errorf("Expected config file to contain '%s'", expected)
		}
	}
}

func TestAuthMethod_IsValid(t *testing.T) {
	testCases := []struct {
		method types.AuthMethod
		valid  bool
	}{
		{types.AuthMethodToken, true},
		{types.AuthMethodUserPass, true},
		{types.AuthMethodLDAP, true},
		{types.AuthMethodOIDC, true},
		{types.AuthMethod("invalid"), false},
		{types.AuthMethod(""), false},
	}

	for _, tc := range testCases {
		if tc.method.IsValid() != tc.valid {
			t.Errorf("Expected %s.IsValid() to be %v, got %v", tc.method, tc.valid, tc.method.IsValid())
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
