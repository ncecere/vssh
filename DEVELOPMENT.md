# vssh Development Guide

This document provides comprehensive information for developers working on the vssh project.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Building the Application](#building-the-application)
- [Testing](#testing)
- [Code Architecture](#code-architecture)
- [Contributing Guidelines](#contributing-guidelines)
- [Release Process](#release-process)

## Development Environment Setup

### Prerequisites

- **Go 1.19 or later**: [Download Go](https://golang.org/dl/)
- **Git**: For version control
- **Make** (optional): For using Makefile commands
- **HashiCorp Vault**: For testing (can use Vault dev server)
- **SSH client**: OpenSSH client for testing

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ncecere/vssh.git
   cd vssh
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Verify setup**:
   ```bash
   go build -o vssh
   ./vssh --help
   ```

### Development Tools

#### Recommended VS Code Extensions

- **Go** (golang.go): Official Go extension
- **Go Test Explorer**: For running tests
- **YAML**: For configuration file editing
- **GitLens**: Enhanced Git capabilities

#### Useful Go Tools

```bash
# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/sast-scan@latest
```

## Project Structure

```
vssh/
├── cmd/                     # CLI commands and entry points
│   ├── init.go             # Initialize configuration command
│   └── root.go             # Root command and main CLI logic
├── internal/               # Private application code
│   ├── auth/               # Vault authentication logic
│   │   └── auth.go         # Authentication methods implementation
│   ├── config/             # Configuration management
│   │   └── config.go       # Config loading, validation, defaults
│   ├── ssh/                # SSH operations
│   │   ├── client.go       # SSH client wrapper and execution
│   │   └── signing.go      # SSH key signing and certificate management
│   ├── utils/              # Utility functions
│   │   └── logger.go       # Logging configuration
│   └── vault/              # Vault client wrapper
│       └── client.go       # Vault API client and token management
├── pkg/                    # Public API (importable by other projects)
│   └── types/              # Public type definitions
│       └── config.go       # Configuration structures
├── tests/                  # Test files
│   └── config/             # Configuration tests
│       └── config_test.go  # Unit tests for config package
├── docs/                   # Documentation (optional)
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point
├── README.md               # User documentation
├── DEVELOPMENT.md          # This file
├── PLANNING.md             # Project planning and architecture
├── TASK.md                 # Development task tracking
└── prd.md                  # Product requirements document
```

### Package Responsibilities

#### `cmd/`
- **Purpose**: CLI command definitions and argument parsing
- **Key Files**:
  - `root.go`: Main command logic, orchestrates the entire flow
  - `init.go`: Configuration initialization command

#### `internal/auth/`
- **Purpose**: Vault authentication methods
- **Key Components**:
  - Token authentication
  - Username/Password authentication
  - LDAP authentication
  - OIDC authentication
  - Interactive credential prompting

#### `internal/config/`
- **Purpose**: Configuration management
- **Key Components**:
  - Configuration file loading and validation
  - Default value setting
  - Environment variable integration
  - XDG configuration directory support

#### `internal/ssh/`
- **Purpose**: SSH operations and certificate management
- **Key Components**:
  - SSH target parsing (`user@hostname`)
  - Certificate signing with Vault
  - Certificate validation and caching
  - SSH client execution

#### `internal/vault/`
- **Purpose**: Vault API client wrapper
- **Key Components**:
  - Vault client initialization
  - Token management and validation
  - API request handling

#### `pkg/types/`
- **Purpose**: Public type definitions
- **Key Components**:
  - Configuration structures
  - Authentication method enums
  - Public interfaces

## Building the Application

### Basic Build

```bash
# Build for current platform
go build -o vssh

# Build with version information
go build -ldflags "-X main.version=v1.0.0" -o vssh

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o vssh-linux-amd64
```

### Cross-Platform Builds

```bash
# Build for multiple platforms
make build-all

# Or manually:
GOOS=linux GOARCH=amd64 go build -o dist/vssh-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o dist/vssh-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o dist/vssh-windows-amd64.exe
```

### Development Build

```bash
# Build with debug information
go build -gcflags="all=-N -l" -o vssh-debug

# Build and run
go run main.go --help
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config/

# Run specific test
go test -run TestLoadConfig ./internal/config/
```

### Test Structure

Tests are organized alongside the code they test:

```
internal/config/
├── config.go
└── config_test.go
```

### Writing Tests

Example test structure:

```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name    string
        setup   func() error
        want    *types.Config
        wantErr bool
    }{
        {
            name: "valid config",
            setup: func() error {
                // Setup test configuration
                return nil
            },
            want: &types.Config{
                // Expected configuration
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                if err := tt.setup(); err != nil {
                    t.Fatalf("setup failed: %v", err)
                }
            }

            got, err := LoadConfig()
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Testing

For integration tests with Vault:

```bash
# Start Vault dev server
vault server -dev -dev-root-token-id=root

# Run integration tests
VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root go test -tags=integration ./...
```

## Code Architecture

### Design Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Injection**: Dependencies are injected rather than created internally
3. **Error Handling**: Comprehensive error handling with context
4. **Logging**: Structured logging throughout the application
5. **Configuration**: Centralized configuration management

### Key Patterns

#### Configuration Pattern

```go
// Load configuration once at startup
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

// Pass configuration to components that need it
vaultClient, err := vault.NewClient(&cfg.Vault)
sshClient := ssh.NewClient(cfg, logger)
```

#### Authentication Pattern

```go
// Create authenticator with dependencies
authenticator := auth.NewAuthenticator(vaultClient, &cfg.Vault, logger)

// Ensure authentication (handles all auth methods)
if err := authenticator.EnsureAuthenticated(); err != nil {
    log.Fatal(err)
}
```

#### SSH Operations Pattern

```go
// Parse SSH target
target, err := ssh.ParseSSHTarget(args[0])

// Ensure certificate exists and is valid
signer := ssh.NewSigner(vaultClient, cfg, logger)
certPath, err := signer.EnsureSSHCertificate(target.Username)

// Execute SSH connection
sshClient := ssh.NewClient(cfg, logger)
err = sshClient.Connect(target, certPath, options, command)
```

### Error Handling

Use wrapped errors for better context:

```go
if err != nil {
    return fmt.Errorf("failed to load configuration: %w", err)
}
```

### Logging

Use structured logging with logrus:

```go
logger.WithFields(logrus.Fields{
    "username": username,
    "hostname": hostname,
}).Info("Connecting to server")
```

## Contributing Guidelines

### Code Style

1. **Follow Go conventions**: Use `gofmt`, `goimports`, and `golint`
2. **Naming**: Use descriptive names for functions and variables
3. **Comments**: Document public functions and complex logic
4. **Error messages**: Provide clear, actionable error messages

### Commit Messages

Use conventional commit format:

```
feat: add OIDC authentication support
fix: resolve certificate validation issue
docs: update configuration examples
test: add unit tests for SSH client
```

### Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Make changes** with tests
4. **Run tests**: `go test ./...`
5. **Run linting**: `golangci-lint run`
6. **Commit changes** with clear messages
7. **Push to fork**: `git push origin feature/new-feature`
8. **Create Pull Request**

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] Error handling is comprehensive
- [ ] Logging is appropriate
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed

## Release Process

### Version Management

Use semantic versioning (semver):

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

### Release Steps

1. **Update version** in relevant files
2. **Update CHANGELOG.md**
3. **Create release branch**: `git checkout -b release/v1.2.0`
4. **Run full test suite**
5. **Build release binaries**
6. **Create Git tag**: `git tag v1.2.0`
7. **Push tag**: `git push origin v1.2.0`
8. **Create GitHub release**
9. **Update documentation**

### Build Scripts

Create a `Makefile` for common tasks:

```makefile
.PHONY: build test clean release

build:
	go build -o vssh

test:
	go test ./...

clean:
	rm -f vssh vssh-*

release:
	GOOS=linux GOARCH=amd64 go build -o dist/vssh-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/vssh-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o dist/vssh-windows-amd64.exe
```

## Debugging

### Debug Build

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o vssh-debug

# Use with delve debugger
dlv exec ./vssh-debug -- --debug user@server.com
```

### Logging Levels

```bash
# Enable debug logging
./vssh --debug user@server.com

# Enable verbose logging
./vssh --verbose user@server.com
```

### Common Debug Scenarios

1. **Authentication Issues**: Check Vault connectivity and credentials
2. **Certificate Problems**: Verify SSH key pairs and Vault roles
3. **SSH Connection Failures**: Test manual SSH with generated certificates
4. **Configuration Errors**: Validate YAML syntax and required fields

## Performance Considerations

### Optimization Areas

1. **Certificate Caching**: Reuse valid certificates
2. **Token Caching**: Avoid unnecessary Vault authentication
3. **Concurrent Operations**: Consider parallel processing for multiple connections
4. **Memory Usage**: Minimize memory allocation in hot paths

### Profiling

```bash
# CPU profiling
go build -o vssh
./vssh -cpuprofile=cpu.prof user@server.com
go tool pprof cpu.prof

# Memory profiling
go build -o vssh
./vssh -memprofile=mem.prof user@server.com
go tool pprof mem.prof
```

## Security Considerations

### Development Security

1. **Secrets Management**: Never commit secrets to version control
2. **Dependencies**: Regularly update dependencies for security patches
3. **Input Validation**: Validate all user inputs
4. **File Permissions**: Use appropriate file permissions for sensitive files

### Security Testing

```bash
# Run security scanner
gosec ./...

# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

## Troubleshooting Development Issues

### Common Issues

1. **Module Issues**: Run `go mod tidy` to clean up dependencies
2. **Build Failures**: Check Go version and environment variables
3. **Test Failures**: Ensure test dependencies are available
4. **Import Errors**: Verify module paths and Go workspace setup

### Getting Help

1. **Documentation**: Check README.md and code comments
2. **Issues**: Search existing GitHub issues
3. **Community**: Join project discussions
4. **Debugging**: Use debug logging and Go debugging tools

This development guide should provide everything needed to contribute effectively to the vssh project. For questions or improvements to this guide, please open an issue or submit a pull request.
