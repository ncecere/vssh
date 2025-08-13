package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vssh/internal/utils"
	"vssh/pkg/types"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

// Client wraps the Vault API client with additional functionality
type Client struct {
	client *api.Client
	config *types.VaultConfig
	logger *logrus.Logger
}

// NewClient creates a new Vault client
func NewClient(config *types.VaultConfig) (*Client, error) {
	// Create Vault client configuration
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = config.Address

	// Create the client
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// Set namespace if configured
	if config.Namespace != "" {
		client.SetNamespace(config.Namespace)
	}

	return &Client{
		client: client,
		config: config,
		logger: utils.GetLogger(),
	}, nil
}

// IsTokenValid checks if the current token is valid and not expired
func (c *Client) IsTokenValid() bool {
	// Get current token
	token := c.client.Token()
	if token == "" {
		c.logger.Debug("No token found")
		return false
	}

	// Check token validity by looking up self
	secret, err := c.client.Auth().Token().LookupSelf()
	if err != nil {
		c.logger.Debugf("Token lookup failed: %v", err)
		return false
	}

	// Check if token is renewable and get TTL
	if secret.Data == nil {
		c.logger.Debug("Token lookup returned no data")
		return false
	}

	// Get TTL from response
	ttlInterface, exists := secret.Data["ttl"]
	if !exists {
		c.logger.Debug("Token TTL not found in response")
		return false
	}

	// Convert TTL to duration
	var ttl time.Duration
	switch v := ttlInterface.(type) {
	case int:
		ttl = time.Duration(v) * time.Second
	case int64:
		ttl = time.Duration(v) * time.Second
	case float64:
		ttl = time.Duration(v) * time.Second
	case json.Number:
		if ttlInt, err := v.Int64(); err == nil {
			ttl = time.Duration(ttlInt) * time.Second
		} else {
			c.logger.Debugf("Failed to parse json.Number TTL: %v", err)
			return false
		}
	default:
		c.logger.Debugf("Unexpected TTL type: %T", v)
		return false
	}

	// Consider token valid if it has more than 5 minutes remaining
	minValidTime := 5 * time.Minute
	if ttl < minValidTime {
		c.logger.Debugf("Token TTL too low: %v", ttl)
		return false
	}

	c.logger.Debugf("Token is valid with TTL: %v", ttl)
	return true
}

// LoadTokenFromFile loads a token from the configured token file
func (c *Client) LoadTokenFromFile() error {
	tokenPath := c.config.Token.TokenPath
	if tokenPath == "" {
		return fmt.Errorf("token path not configured")
	}

	// Expand tilde in path
	if tokenPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		tokenPath = home + tokenPath[1:]
	}

	// Read token from file
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		return fmt.Errorf("error reading token file %s: %w", tokenPath, err)
	}

	token := string(tokenBytes)
	// Remove any trailing newlines
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("token file is empty")
	}

	// Set token on client
	c.client.SetToken(token)
	c.logger.Debugf("Loaded token from %s", tokenPath)

	return nil
}

// SaveTokenToFile saves the current token to the configured token file
func (c *Client) SaveTokenToFile() error {
	token := c.client.Token()
	if token == "" {
		return fmt.Errorf("no token to save")
	}

	tokenPath := c.config.Token.TokenPath
	if tokenPath == "" {
		return fmt.Errorf("token path not configured")
	}

	// Expand tilde in path
	if tokenPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		tokenPath = home + tokenPath[1:]
	}

	// Ensure directory exists
	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error creating token directory: %w", err)
	}

	// Write token to file with secure permissions
	err := os.WriteFile(tokenPath, []byte(token), 0600)
	if err != nil {
		return fmt.Errorf("error writing token file: %w", err)
	}

	c.logger.Debugf("Saved token to %s", tokenPath)
	return nil
}

// GetClient returns the underlying Vault API client
func (c *Client) GetClient() *api.Client {
	return c.client
}

// SetToken sets the token on the client
func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}
