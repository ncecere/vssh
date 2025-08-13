# vssh Configuration Reference

This document provides a comprehensive reference for configuring vssh, including all available options, examples, and best practices.

## Table of Contents

- [Configuration File Location](#configuration-file-location)
- [Configuration Structure](#configuration-structure)
- [Vault Configuration](#vault-configuration)
- [SSH Configuration](#ssh-configuration)
- [User Configuration](#user-configuration)
- [Authentication Methods](#authentication-methods)
- [Environment Variables](#environment-variables)
- [Configuration Examples](#configuration-examples)
- [Troubleshooting](#troubleshooting)

## Configuration File Location

vssh follows XDG Base Directory Specification for configuration files:

### Default Location
```
~/.config/vssh/config.yaml
```

### Custom Location
You can specify a custom configuration file using the `--config` flag:
```bash
vssh --config /path/to/custom/config.yaml user@server.com
```

### Configuration File Creation
Initialize a default configuration file:
```bash
vssh init
```

## Configuration Structure

The configuration file uses YAML format with the following top-level sections:

```yaml
vault:      # Vault server and authentication settings
ssh:        # SSH-related configuration
users:      # Per-user SSH key and role configuration
debug:      # Global debug logging setting
```

## Vault Configuration

The `vault` section configures connection and authentication to your HashiCorp Vault server.

### Basic Vault Configuration

```yaml
vault:
  address: "https://vault.example.com:8200"
  auth_method: "token"
  namespace: "my-namespace"  # Optional
```

### Configuration Options

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `address` | string | **Yes** | Vault server URL including protocol and port | - |
| `auth_method` | string | **Yes** | Authentication method: `token`, `userpass`, `ldap`, `oidc` | `token` |
| `namespace` | string | No | Vault namespace (Vault Enterprise feature) | - |

### Vault Address Examples

```yaml
# HTTPS with custom port
address: "https://vault.company.com:8200"

# HTTP (not recommended for production)
address: "http://localhost:8200"

# With path prefix
address: "https://vault.company.com/vault"
```

## SSH Configuration

The `ssh` section configures SSH key management and certificate settings.

### SSH Configuration Options

```yaml
ssh:
  key_directory: "~/.ssh"
  certificate_ttl: "4h"
  signing_engine: "ssh-client-signer"
```

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `key_directory` | string | **Yes** | Directory containing SSH keys | `~/.ssh` |
| `certificate_ttl` | duration | **Yes** | Certificate validity period | `4h` |
| `signing_engine` | string | **Yes** | Vault SSH secrets engine mount path | `ssh-client-signer` |

### Certificate TTL Examples

```yaml
# Various TTL formats
certificate_ttl: "4h"      # 4 hours
certificate_ttl: "30m"     # 30 minutes
certificate_ttl: "1h30m"   # 1 hour 30 minutes
certificate_ttl: "24h"     # 24 hours
certificate_ttl: "7d"      # 7 days (if Vault policy allows)
```

### Key Directory Examples

```yaml
# Default SSH directory
key_directory: "~/.ssh"

# Custom directory
key_directory: "/home/user/keys"

# Absolute path
key_directory: "/opt/ssh-keys"
```

## User Configuration

The `users` section allows per-user customization of SSH keys and Vault roles.

### User Configuration Structure

```yaml
users:
  username1:
    private_key: "~/.ssh/username1_rsa"
    vault_role: "custom-role"
  username2:
    private_key: "~/.ssh/username2_rsa"
    vault_role: "username2-special"
```

### User Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `private_key` | string | **Yes** | Path to user's private SSH key |
| `vault_role` | string | No | Custom Vault role (defaults to username) |

### User Configuration Examples

#### Basic Multi-User Setup
```yaml
users:
  alice:
    private_key: "~/.ssh/alice_rsa"
  bob:
    private_key: "~/.ssh/bob_rsa"
  charlie:
    private_key: "~/.ssh/charlie_rsa"
```

#### Custom Vault Roles
```yaml
users:
  admin:
    private_key: "~/.ssh/admin_rsa"
    vault_role: "ssh-admin-role"
  developer:
    private_key: "~/.ssh/dev_rsa"
    vault_role: "ssh-developer-role"
  readonly:
    private_key: "~/.ssh/readonly_rsa"
    vault_role: "ssh-readonly-role"
```

#### Mixed Key Types
```yaml
users:
  user1:
    private_key: "~/.ssh/id_rsa"        # RSA key
  user2:
    private_key: "~/.ssh/id_ed25519"    # Ed25519 key
  user3:
    private_key: "~/.ssh/id_ecdsa"      # ECDSA key
```

## Authentication Methods

vssh supports four authentication methods for connecting to Vault.

### Token Authentication (Default)

Uses a Vault token stored in a file.

```yaml
vault:
  auth_method: "token"
  token:
    token_path: "~/.vault-token"
```

#### Token Configuration Options

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `token_path` | string | No | Path to Vault token file | `~/.vault-token` |

#### Token Authentication Examples

```yaml
# Default token location
vault:
  auth_method: "token"
  token:
    token_path: "~/.vault-token"

# Custom token location
vault:
  auth_method: "token"
  token:
    token_path: "/opt/vault/token"

# Environment-based token path
vault:
  auth_method: "token"
  token:
    token_path: "${HOME}/.config/vault/token"
```

### Username/Password Authentication

Uses username and password for authentication.

```yaml
vault:
  auth_method: "userpass"
  userpass:
    username: "your-username"
    mount: "userpass"
```

#### UserPass Configuration Options

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `username` | string | No | Username (prompted if not provided) | - |
| `mount` | string | No | Auth method mount path | `userpass` |

#### UserPass Examples

```yaml
# Basic userpass with username
vault:
  auth_method: "userpass"
  userpass:
    username: "john.doe"

# Custom mount path
vault:
  auth_method: "userpass"
  userpass:
    username: "john.doe"
    mount: "corporate-userpass"

# Username prompted at runtime
vault:
  auth_method: "userpass"
  userpass:
    mount: "userpass"
```

### LDAP Authentication

Uses LDAP credentials for authentication.

```yaml
vault:
  auth_method: "ldap"
  ldap:
    username: "your-username"
    mount: "ldap"
```

#### LDAP Configuration Options

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `username` | string | No | LDAP username (prompted if not provided) | - |
| `mount` | string | No | Auth method mount path | `ldap` |

#### LDAP Examples

```yaml
# Basic LDAP with username
vault:
  auth_method: "ldap"
  ldap:
    username: "john.doe"

# Custom mount path
vault:
  auth_method: "ldap"
  ldap:
    username: "john.doe"
    mount: "corporate-ldap"

# Username prompted at runtime
vault:
  auth_method: "ldap"
  ldap:
    mount: "ldap"
```

### OIDC Authentication

Uses OpenID Connect for authentication.

```yaml
vault:
  auth_method: "oidc"
  oidc:
    role: "your-oidc-role"
    mount: "oidc"
```

#### OIDC Configuration Options

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `role` | string | **Yes** | OIDC role name | - |
| `mount` | string | No | Auth method mount path | `oidc` |

#### OIDC Examples

```yaml
# Basic OIDC
vault:
  auth_method: "oidc"
  oidc:
    role: "engineering"

# Custom mount path
vault:
  auth_method: "oidc"
  oidc:
    role: "engineering"
    mount: "company-oidc"

# Multiple roles (use different configs)
vault:
  auth_method: "oidc"
  oidc:
    role: "admin"
    mount: "oidc"
```

## Environment Variables

vssh respects several environment variables that can override configuration settings.

### Supported Environment Variables

| Variable | Description | Configuration Override |
|----------|-------------|----------------------|
| `VAULT_ADDR` | Vault server address | `vault.address` |
| `VAULT_TOKEN` | Vault token | Used for token auth |
| `VAULT_NAMESPACE` | Vault namespace | `vault.namespace` |
| `USER` | Current username | Used as fallback username |

### Environment Variable Examples

```bash
# Override Vault address
export VAULT_ADDR=https://vault.company.com:8200
vssh user@server.com

# Use environment token
export VAULT_TOKEN=hvs.CAESIJ...
vssh user@server.com

# Set namespace
export VAULT_NAMESPACE=engineering
vssh user@server.com
```

## Configuration Examples

### Complete Production Configuration

```yaml
# Production configuration for enterprise environment
vault:
  address: "https://vault.company.com:8200"
  auth_method: "ldap"
  namespace: "engineering"
  
  ldap:
    mount: "corporate-ldap"

ssh:
  key_directory: "~/.ssh"
  certificate_ttl: "8h"
  signing_engine: "ssh-client-signer"

users:
  # Development team
  alice:
    private_key: "~/.ssh/alice_rsa"
    vault_role: "ssh-developer"
  
  bob:
    private_key: "~/.ssh/bob_rsa"
    vault_role: "ssh-developer"
  
  # Operations team
  charlie:
    private_key: "~/.ssh/charlie_rsa"
    vault_role: "ssh-admin"
  
  # Service accounts
  deploy:
    private_key: "~/.ssh/deploy_rsa"
    vault_role: "ssh-deployment"

debug: false
```

### Development Configuration

```yaml
# Development configuration for local testing
vault:
  address: "http://localhost:8200"
  auth_method: "token"
  
  token:
    token_path: "~/.vault-token"

ssh:
  key_directory: "~/.ssh"
  certificate_ttl: "1h"
  signing_engine: "ssh"

users:
  dev:
    private_key: "~/.ssh/id_rsa"

debug: true
```

### Multi-Environment Configuration

```yaml
# Configuration supporting multiple environments
vault:
  address: "https://vault-prod.company.com:8200"
  auth_method: "oidc"
  namespace: "production"
  
  oidc:
    role: "engineering"
    mount: "company-oidc"

ssh:
  key_directory: "~/.ssh"
  certificate_ttl: "4h"
  signing_engine: "ssh-client-signer"

users:
  # Production users
  prod-admin:
    private_key: "~/.ssh/prod_admin_rsa"
    vault_role: "ssh-prod-admin"
  
  prod-deploy:
    private_key: "~/.ssh/prod_deploy_rsa"
    vault_role: "ssh-prod-deploy"
  
  # Staging users
  staging-admin:
    private_key: "~/.ssh/staging_admin_rsa"
    vault_role: "ssh-staging-admin"

debug: false
```

### Minimal Configuration

```yaml
# Minimal configuration with defaults
vault:
  address: "https://vault.example.com:8200"
  auth_method: "token"

ssh:
  key_directory: "~/.ssh"
  certificate_ttl: "4h"
  signing_engine: "ssh-client-signer"
```

## Troubleshooting

### Configuration Validation

vssh validates configuration on startup. Common validation errors:

#### Missing Required Fields
```
Error: vault.address is required
```
**Solution**: Add the Vault server address to your configuration.

#### Invalid Authentication Method
```
Error: invalid auth method: invalid. Supported methods: token, userpass, ldap, oidc
```
**Solution**: Use a supported authentication method.

#### Invalid TTL Format
```
Error: time: invalid duration "invalid"
```
**Solution**: Use valid duration format (e.g., "4h", "30m", "1h30m").

### Configuration File Issues

#### File Not Found
```
Error: Config File "config" Not Found in "[/home/user/.config/vssh]"
```
**Solution**: Run `vssh init` to create a default configuration file.

#### YAML Syntax Errors
```
Error: yaml: line 5: mapping values are not allowed in this context
```
**Solution**: Check YAML syntax, ensure proper indentation and structure.

#### Permission Errors
```
Error: open /home/user/.config/vssh/config.yaml: permission denied
```
**Solution**: Check file permissions and ownership.

### Authentication Issues

#### Token Authentication
```
Error: Error making API request. Code: 403. Errors: * permission denied
```
**Solutions**:
- Check token validity: `vault token lookup`
- Verify token has required permissions
- Renew token if expired

#### LDAP Authentication
```
Error: Error making API request. Code: 400. Errors: * invalid username or password
```
**Solutions**:
- Verify LDAP credentials
- Check LDAP mount path
- Ensure LDAP auth method is enabled in Vault

#### OIDC Authentication
```
Error: Error making API request. Code: 400. Errors: * role "invalid-role" not found
```
**Solutions**:
- Verify OIDC role exists in Vault
- Check OIDC mount path
- Ensure OIDC auth method is configured

### SSH Configuration Issues

#### Key Directory Not Found
```
Error: stat /home/user/.ssh: no such file or directory
```
**Solution**: Create SSH directory or update `key_directory` path.

#### Private Key Not Found
```
Error: private key not found: /home/user/.ssh/id_rsa
```
**Solutions**:
- Generate SSH key pair: `ssh-keygen -t rsa`
- Update private key path in configuration
- Check file permissions

#### Certificate TTL Too Long
```
Error: Error making API request. Code: 400. Errors: * ttl is larger than max_ttl
```
**Solution**: Reduce `certificate_ttl` value or increase Vault role's `max_ttl`.

### Debug Configuration

Enable debug logging to troubleshoot configuration issues:

```yaml
debug: true
```

Or use the command-line flag:
```bash
vssh --debug user@server.com
```

### Configuration Testing

Test configuration without making SSH connections:

```bash
# Test configuration loading
vssh --help

# Test Vault authentication
vault auth -method=ldap username=your-username

# Test SSH key access
ls -la ~/.ssh/
```

This configuration reference should help you set up vssh for any environment. For additional help, consult the main README.md or open an issue on GitHub.
