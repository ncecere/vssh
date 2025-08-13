package types

import "time"

// Config represents the main configuration structure
type Config struct {
	Vault VaultConfig `mapstructure:"vault" yaml:"vault"`
	SSH   SSHConfig   `mapstructure:"ssh" yaml:"ssh"`
	Users UserConfigs `mapstructure:"users" yaml:"users"`
	Debug bool        `mapstructure:"debug" yaml:"debug"`
}

// VaultConfig contains Vault server configuration
type VaultConfig struct {
	Address    string `mapstructure:"address" yaml:"address"`
	Role       string `mapstructure:"role" yaml:"role"`
	AuthMethod string `mapstructure:"auth_method" yaml:"auth_method"`
	Namespace  string `mapstructure:"namespace" yaml:"namespace,omitempty"`

	// Auth method specific configurations
	Token    TokenConfig    `mapstructure:"token" yaml:"token,omitempty"`
	UserPass UserPassConfig `mapstructure:"userpass" yaml:"userpass,omitempty"`
	LDAP     LDAPConfig     `mapstructure:"ldap" yaml:"ldap,omitempty"`
	OIDC     OIDCConfig     `mapstructure:"oidc" yaml:"oidc,omitempty"`
}

// TokenConfig for token-based authentication
type TokenConfig struct {
	TokenPath string `mapstructure:"token_path" yaml:"token_path,omitempty"`
}

// UserPassConfig for username/password authentication
type UserPassConfig struct {
	Username string `mapstructure:"username" yaml:"username"`
	Mount    string `mapstructure:"mount" yaml:"mount,omitempty"`
}

// LDAPConfig for LDAP authentication
type LDAPConfig struct {
	Username string `mapstructure:"username" yaml:"username"`
	Mount    string `mapstructure:"mount" yaml:"mount,omitempty"`
}

// OIDCConfig for OIDC authentication
type OIDCConfig struct {
	Role  string `mapstructure:"role" yaml:"role"`
	Mount string `mapstructure:"mount" yaml:"mount,omitempty"`
}

// SSHConfig contains SSH-related configuration
type SSHConfig struct {
	KeyDirectory   string        `mapstructure:"key_directory" yaml:"key_directory"`
	CertificateTTL time.Duration `mapstructure:"certificate_ttl" yaml:"certificate_ttl"`
	SigningEngine  string        `mapstructure:"signing_engine" yaml:"signing_engine"`
}

// UserConfig represents per-user configuration
type UserConfig struct {
	PrivateKey string `mapstructure:"private_key" yaml:"private_key"`
	VaultRole  string `mapstructure:"vault_role" yaml:"vault_role,omitempty"`
}

// UserConfigs is a map of username to user configuration
type UserConfigs map[string]UserConfig

// AuthMethod represents supported authentication methods
type AuthMethod string

const (
	AuthMethodToken    AuthMethod = "token"
	AuthMethodUserPass AuthMethod = "userpass"
	AuthMethodLDAP     AuthMethod = "ldap"
	AuthMethodOIDC     AuthMethod = "oidc"
)

// IsValid checks if the auth method is supported
func (a AuthMethod) IsValid() bool {
	switch a {
	case AuthMethodToken, AuthMethodUserPass, AuthMethodLDAP, AuthMethodOIDC:
		return true
	default:
		return false
	}
}
