# vssh - Vault SSH Certificate Tool

A Go CLI tool that signs SSH keys with HashiCorp Vault and uses signed certificates for SSH authentication. The tool acts as a wrapper around SSH, providing seamless certificate-based authentication through Vault.

## Features

- **Multiple Authentication Methods**: Support for Token, Username/Password, LDAP, and OIDC authentication
- **Username-based Roles**: Automatically uses SSH username as Vault role (matches `vault write ssh-client-signer/sign/username` pattern)
- **Multi-User Support**: Different users can have different SSH keys and Vault roles
- **Automatic Certificate Management**: Validates and renews certificates automatically
- **SSH Compatibility**: Works like regular SSH with additional Vault integration
- **Certificate Caching**: Reuses valid certificates until expiration
- **Comprehensive Logging**: Debug and verbose logging for troubleshooting

## Installation

vssh provides multiple installation methods to suit different preferences and environments.

### Quick Install (Recommended)

#### Linux and macOS
```bash
curl -fsSL https://raw.githubusercontent.com/ncecere/vssh/main/install.sh | bash
```

#### Windows (PowerShell)
```powershell
iwr -useb https://raw.githubusercontent.com/ncecere/vssh/main/install.ps1 | iex
```

### Manual Download

Download the appropriate binary for your platform from the [latest release](https://github.com/ncecere/vssh/releases/latest):

#### Linux
```bash
# AMD64 (Intel/AMD 64-bit)
curl -L https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-linux-amd64 -o vssh
chmod +x vssh
sudo mv vssh /usr/local/bin/

# ARM64 (ARM 64-bit)
curl -L https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-linux-arm64 -o vssh
chmod +x vssh
sudo mv vssh /usr/local/bin/
```

#### macOS
```bash
# Intel Mac
curl -L https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-darwin-amd64 -o vssh
chmod +x vssh
sudo mv vssh /usr/local/bin/

# Apple Silicon (M1/M2)
curl -L https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-darwin-arm64 -o vssh
chmod +x vssh
sudo mv vssh /usr/local/bin/
```

#### Windows
```powershell
# AMD64 (Intel/AMD 64-bit)
Invoke-WebRequest -Uri "https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-windows-amd64.exe" -OutFile "vssh.exe"

# ARM64 (ARM 64-bit)
Invoke-WebRequest -Uri "https://github.com/ncecere/vssh/releases/latest/download/vssh-v0.1.0-windows-arm64.exe" -OutFile "vssh.exe"

# Add to PATH (optional)
$env:PATH += ";$(Get-Location)"
```

### Go Install

If you have Go installed, you can install vssh directly:

```bash
# Install latest version
go install github.com/ncecere/vssh@latest

# Install specific version
go install github.com/ncecere/vssh@v0.1.0
```

### GitHub CLI

If you have GitHub CLI installed:

```bash
# Download latest release
gh release download -R ncecere/vssh

# Download specific version
gh release download v0.1.0 -R ncecere/vssh
```

### Build from Source

For developers who want to build from source:

```bash
git clone https://github.com/ncecere/vssh.git
cd vssh
go mod download
go build -o vssh
sudo mv vssh /usr/local/bin/
```

### Installation Verification

After installation, verify vssh is working:

```bash
# Check version
vssh version

# Show help
vssh --help

# Initialize configuration
vssh init
```

### Prerequisites

- Access to a HashiCorp Vault server with SSH secrets engine enabled
- SSH client installed on your system
- SSH key pair generated (`ssh-keygen -t rsa` or `ssh-keygen -t ed25519`)

### Upgrading

To upgrade to a newer version:

#### Using Install Scripts
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/ncecere/vssh/main/install.sh | bash

# Windows
iwr -useb https://raw.githubusercontent.com/ncecere/vssh/main/install.ps1 | iex
```

#### Using Go Install
```bash
go install github.com/ncecere/vssh@latest
```

#### Manual Upgrade
Download the new binary and replace the existing one following the manual installation steps above.

### Uninstallation

#### Linux/macOS
```bash
# Remove binary
sudo rm /usr/local/bin/vssh

# Remove configuration (optional)
rm -rf ~/.config/vssh

# Remove cached certificates (optional)
rm -f ~/.ssh/vault_signed_*.pub
```

#### Windows
```powershell
# Remove binary (adjust path as needed)
Remove-Item "$env:LOCALAPPDATA\vssh\bin\vssh.exe"

# Remove from PATH if added manually
# Edit environment variables through System Properties

# Remove configuration (optional)
Remove-Item -Recurse "$env:USERPROFILE\.config\vssh"

# Remove cached certificates (optional)
Remove-Item "$env:USERPROFILE\.ssh\vault_signed_*.pub"
```

## Quick Start

1. **Initialize configuration**:
   ```bash
   vssh init
   ```

2. **Edit the configuration file** at `~/.config/vssh/config.yaml`:
   ```yaml
   vault:
     address: "https://your-vault-server.com"
     auth_method: "ldap"  # or token, userpass, oidc
   
   ssh:
     key_directory: "~/.ssh"
     certificate_ttl: "4h"
     signing_engine: "ssh-client-signer"
   ```

3. **Connect to a server**:
   ```bash
   vssh user@server.com
   ```

## Command Line Reference

### Basic Usage

```bash
vssh [flags] [user@]hostname [command]
```

### Global Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--config` | | Custom config file path | `--config /path/to/config.yaml` |
| `--verbose` | `-v` | Enable verbose output | `vssh -v user@server.com` |
| `--debug` | `-d` | Enable debug output | `vssh --debug user@server.com` |
| `--help` | `-h` | Show help information | `vssh --help` |

### SSH-Compatible Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--port` | `-p` | Port to connect to | `vssh -p 2222 user@server.com` |
| `--identity` | `-i` | Identity (private key) file | `vssh -i ~/.ssh/custom_key user@server.com` |
| `--ipv4` | `-4` | Force IPv4 addresses only | `vssh -4 user@server.com` |
| `--ipv6` | `-6` | Force IPv6 addresses only | `vssh -6 user@server.com` |

### Commands

#### Initialize Configuration
```bash
vssh init                    # Create default config file
vssh init --help             # Show init command help
```

### Usage Examples

#### Basic Connection
```bash
# Connect to server with username
vssh user@server.com

# Connect using current system username
vssh server.com

# Run a single command
vssh user@server.com "ls -la /home"

# Interactive session with debug logging
vssh --debug user@server.com
```

#### Advanced Usage
```bash
# Custom port
vssh -p 2222 user@server.com

# Custom identity file
vssh -i ~/.ssh/custom_key user@server.com

# Force IPv4
vssh -4 user@server.com

# Combine multiple options
vssh --debug -p 2222 -4 user@server.com "uptime"
```

## Configuration

The configuration file is located at `~/.config/vssh/config.yaml`. You can specify a custom location with the `--config` flag.

### Complete Configuration Example

```yaml
# Vault server configuration
vault:
  address: "https://vault.example.com:8200"
  auth_method: "ldap"                    # token, userpass, ldap, oidc
  namespace: "my-namespace"              # Optional: Vault namespace
  
  # Token authentication
  token:
    token_path: "~/.vault-token"         # Path to Vault token file
  
  # Username/Password authentication
  userpass:
    username: "your-username"            # Optional: can be prompted
    mount: "userpass"                    # Optional: auth mount path
  
  # LDAP authentication
  ldap:
    username: "your-username"            # Optional: can be prompted
    mount: "ldap"                        # Optional: auth mount path
  
  # OIDC authentication
  oidc:
    role: "your-oidc-role"               # Required for OIDC
    mount: "oidc"                        # Optional: auth mount path

# SSH configuration
ssh:
  key_directory: "~/.ssh"                # Directory containing SSH keys
  certificate_ttl: "4h"                  # Certificate validity period
  signing_engine: "ssh-client-signer"    # Vault SSH secrets engine mount

# Per-user SSH key configuration (optional)
users:
  user1:
    private_key: "~/.ssh/user1_rsa"      # Path to user's private key
    vault_role: "custom-role"            # Optional: override username-based role
  
  user2:
    private_key: "~/.ssh/user2_rsa"
    vault_role: "user2-special-role"

# Enable debug logging globally
debug: false
```

### Configuration Options Reference

#### Vault Section

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `address` | string | Yes | Vault server URL | - |
| `auth_method` | string | Yes | Authentication method | `token` |
| `namespace` | string | No | Vault namespace | - |

#### Authentication Methods

**Token Authentication** (default)
```yaml
vault:
  auth_method: "token"
  token:
    token_path: "~/.vault-token"         # Path to token file
```

**Username/Password Authentication**
```yaml
vault:
  auth_method: "userpass"
  userpass:
    username: "your-username"            # Optional: prompted if not provided
    mount: "userpass"                    # Optional: defaults to "userpass"
```

**LDAP Authentication**
```yaml
vault:
  auth_method: "ldap"
  ldap:
    username: "your-username"            # Optional: prompted if not provided
    mount: "ldap"                        # Optional: defaults to "ldap"
```

**OIDC Authentication**
```yaml
vault:
  auth_method: "oidc"
  oidc:
    role: "your-oidc-role"               # Required
    mount: "oidc"                        # Optional: defaults to "oidc"
```

#### SSH Section

| Option | Type | Required | Description | Default |
|--------|------|----------|-------------|---------|
| `key_directory` | string | Yes | SSH keys directory | `~/.ssh` |
| `certificate_ttl` | duration | Yes | Certificate validity | `4h` |
| `signing_engine` | string | Yes | Vault SSH engine mount | `ssh-client-signer` |

#### Users Section

Per-user configuration for custom SSH keys and Vault roles:

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `private_key` | string | Yes | Path to user's private key |
| `vault_role` | string | No | Custom Vault role (defaults to username) |

## How It Works

### Authentication Flow

1. **Load Configuration**: Read settings from `~/.config/vssh/config.yaml`
2. **Parse SSH Target**: Extract username and hostname from command arguments
3. **Vault Authentication**: Authenticate using configured method (token/userpass/ldap/oidc)
4. **Token Caching**: Cache valid tokens for reuse
5. **Certificate Check**: Verify if existing certificate is still valid
6. **Key Signing**: Request new certificate from Vault if needed using username as role
7. **SSH Execution**: Connect using signed certificate and private key

### Certificate Management

- **Naming Convention**: Certificates are named `vault_signed_{username}.pub`
- **Role Mapping**: Username is used as Vault role (e.g., `user1@server.com` → role `user1`)
- **Automatic Renewal**: Certificates are renewed when they expire or have <5 minutes remaining
- **Validation**: Certificates are validated before each use
- **Storage**: Certificates are stored in the configured SSH key directory

### Vault Integration

vssh integrates with Vault's SSH secrets engine:

```bash
# Equivalent Vault CLI command
vault write ssh-client-signer/sign/username \
  public_key=@~/.ssh/id_rsa.pub \
  ttl=4h
```

The tool automatically:
- Uses the SSH username as the Vault role name
- Reads the corresponding public key
- Requests certificate signing
- Saves the signed certificate
- Uses it for SSH authentication

## Troubleshooting

### Common Issues

#### Authentication Failures
```bash
# Check Vault connectivity
curl -k https://your-vault-server.com/v1/sys/health

# Test authentication manually
vault auth -method=ldap username=your-username

# Enable debug logging
vssh --debug user@server.com
```

#### Certificate Issues
```bash
# Check certificate validity
ssh-keygen -L -f ~/.ssh/vault_signed_username.pub

# Force certificate renewal by removing existing certificate
rm ~/.ssh/vault_signed_username.pub
vssh user@server.com
```

#### SSH Connection Issues
```bash
# Test SSH manually with certificate
ssh -o CertificateFile=~/.ssh/vault_signed_username.pub \
    -i ~/.ssh/id_rsa \
    user@server.com

# Check SSH server configuration
# Ensure TrustedUserCAKeys is configured on the server
```

### Debug Logging

Enable debug logging to troubleshoot issues:

```bash
vssh --debug user@server.com
```

Debug output includes:
- Configuration loading
- Vault authentication details
- Certificate validation
- SSH command construction
- Connection attempts

### Environment Variables

vssh respects standard environment variables:

| Variable | Description |
|----------|-------------|
| `VAULT_ADDR` | Vault server address (overrides config) |
| `VAULT_TOKEN` | Vault token (for token auth) |
| `VAULT_NAMESPACE` | Vault namespace |
| `USER` | Current username (fallback if not specified) |

## Security Considerations

- **Token Storage**: Vault tokens are cached with secure file permissions (0600)
- **Private Key Protection**: Private keys remain on the local system and are never sent to Vault
- **Certificate Lifecycle**: Certificates have configurable TTL and are automatically renewed
- **Audit Logging**: All authentication and signing events are logged by Vault
- **Least Privilege**: Each user can only sign certificates for their own username role

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed development instructions.

### Quick Development Setup

```bash
git clone https://github.com/ncecere/vssh.git
cd vssh
go mod download
go build -o vssh
./vssh init
```

## Documentation

- **[README.md](README.md)** - This file, comprehensive user guide
- **[CONFIG.md](CONFIG.md)** - Detailed configuration reference
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Developer setup and build instructions
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Guidelines for contributors
- **[SECURITY.md](SECURITY.md)** - Security policy and vulnerability reporting
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up the development environment
- Code style guidelines
- Testing requirements
- Pull request process

Quick start for contributors:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Security

Security is important to us. If you discover a security vulnerability, please see our [Security Policy](SECURITY.md) for responsible disclosure guidelines.

For general security best practices when using vssh, refer to the Security Considerations section above and the [Security Policy](SECURITY.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs and request features on [GitHub Issues](https://github.com/ncecere/vssh/issues)
- **Security**: Report security vulnerabilities via our [Security Policy](SECURITY.md)
- **Documentation**: Comprehensive documentation available in the repository
- **Community**: Join discussions in [GitHub Discussions](https://github.com/ncecere/vssh/discussions)

## Project Status

- **Current Version**: v0.1.5
- **Status**: Active development
- **Stability**: Beta - suitable for testing and development environments
- **Platform Support**: Linux, macOS, Windows (AMD64, ARM64)

## Roadmap

See [TASK.md](TASK.md) for current development tasks and [CHANGELOG.md](CHANGELOG.md) for version history.

Upcoming features:
- Enhanced error handling and user experience
- Additional authentication methods
- Performance optimizations
- Integration tests

## Acknowledgments

- [HashiCorp Vault](https://www.vaultproject.io/) for the excellent secrets management platform
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management
- [Logrus](https://github.com/sirupsen/logrus) for structured logging
- The Go community for excellent tooling and libraries
