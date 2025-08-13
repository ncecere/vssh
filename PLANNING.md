# vssh - Vault SSH Certificate Tool

## Project Overview
A Go CLI tool that signs SSH keys with HashiCorp Vault and uses signed certificates for SSH authentication. The tool acts as a wrapper around SSH, providing seamless certificate-based authentication through Vault.

## Architecture

### Technology Stack
- **Go** (primary language)
- **Cobra** (CLI framework)
- **Viper** (configuration management)
- **HashiCorp Vault Go Client** (Vault integration)
- **Go SSH library** (SSH operations)
- **Logrus** (structured logging)

### Project Structure
```
vssh/
├── cmd/
│   └── root.go              # Cobra root command
├── internal/
│   ├── auth/                # Vault authentication methods
│   ├── config/              # Configuration management
│   ├── ssh/                 # SSH key operations
│   ├── vault/               # Vault client management
│   └── utils/               # Utilities (logging, files)
├── pkg/
│   └── types/               # Public types
├── tests/                   # Unit tests
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Key Features

### Authentication Methods
- Token authentication
- Username/Password authentication
- LDAP authentication
- OIDC authentication

### Configuration
- XDG standard config location: `~/.config/vssh/config.yaml`
- Per-user SSH key configuration
- Configurable certificate TTL (default: 4 hours)
- Vault server parameters

### SSH Integration
- Multi-user support with separate signed certificates
- Certificate naming: `vault_signed_{username}.pub`
- Automatic certificate validation and renewal
- SSH argument passthrough

### Error Handling
- Fail gracefully with clear error messages
- No fallback to regular SSH (security by design)
- Comprehensive logging for debugging

## Security Considerations
- Secure token caching with proper file permissions
- Private key protection
- Expired certificate cleanup
- Audit logging for authentication events

## Implementation Phases
1. Core Infrastructure (CLI, config, logging)
2. Vault Integration (client, token auth, SSH signing)
3. SSH Operations (key management, certificate handling)
4. Advanced Authentication (userpass, LDAP, OIDC)
5. Polish & Testing (error handling, tests, docs)
