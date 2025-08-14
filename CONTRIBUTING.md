# Contributing to vssh

Thank you for your interest in contributing to vssh! This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Release Process](#release-process)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a feature branch
5. Make your changes
6. Test your changes
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/vssh.git
cd vssh

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

## Making Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-new-auth-method`
- `fix/certificate-validation-bug`
- `docs/update-installation-guide`

### Commit Messages

Follow conventional commit format:
- `feat: add OIDC authentication support`
- `fix: resolve certificate expiration handling`
- `docs: update configuration examples`
- `test: add unit tests for SSH signing`

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./tests/config/

# Run tests with verbose output
go test -v ./...
```

### Writing Tests

- Write unit tests for all new functionality
- Include at least one test for expected behavior
- Include edge case tests
- Include failure case tests
- Place tests in the `tests/` directory mirroring the main structure

### Test Requirements

All pull requests must:
- Include tests for new functionality
- Maintain or improve test coverage
- Pass all existing tests

## Submitting Changes

### Pull Request Process

1. **Update Documentation**: Ensure README.md, CONFIG.md, and other docs are updated
2. **Add Tests**: Include comprehensive tests for your changes
3. **Update Changelog**: Add entry to CHANGELOG.md under "Unreleased"
4. **Check CI**: Ensure all GitHub Actions workflows pass
5. **Request Review**: Submit PR and request review from maintainers

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Changelog updated
```

## Code Style

### Go Style Guidelines

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Follow Go naming conventions
- Use structured logging with logrus

### Documentation Style

- Use clear, concise language
- Include code examples where helpful
- Update all relevant documentation files
- Use proper markdown formatting

### Example Function Documentation

```go
// SignSSHKey signs an SSH public key using Vault's SSH secrets engine.
//
// Args:
//     publicKey (string): The SSH public key to sign
//     username (string): The username for the certificate
//     role (string): The Vault role to use for signing
//
// Returns:
//     string: The signed SSH certificate
//     error: Any error that occurred during signing
func SignSSHKey(publicKey, username, role string) (string, error) {
    // Implementation here
}
```

## Release Process

### Version Management

- Follow semantic versioning (SemVer)
- Update version in relevant files
- Create git tags for releases
- Update CHANGELOG.md with release notes

### Release Checklist

1. Update version numbers
2. Update CHANGELOG.md
3. Run full test suite
4. Create and push git tag
5. GitHub Actions will handle the rest

## Getting Help

- Check existing issues and discussions
- Read the documentation thoroughly
- Ask questions in GitHub issues
- Join discussions in GitHub Discussions

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Follow GitHub's community guidelines

## License

By contributing to vssh, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions about contributing, please:
1. Check this document first
2. Search existing issues
3. Create a new issue with the "question" label
4. Tag maintainers if needed

Thank you for contributing to vssh!
