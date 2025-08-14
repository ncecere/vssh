# Security Policy

## Supported Versions

We actively support the following versions of vssh with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in vssh, please report it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by:

1. **Email**: Send details to the maintainers (create a GitHub issue requesting contact information)
2. **GitHub Security Advisories**: Use GitHub's private vulnerability reporting feature
3. **Direct Message**: Contact maintainers directly through GitHub

### What to Include

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Impact**: What could an attacker accomplish?
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Environment**: Version of vssh, operating system, Go version
- **Proof of Concept**: If applicable, include a minimal example
- **Suggested Fix**: If you have ideas for how to fix it

### Response Timeline

- **Acknowledgment**: We will acknowledge receipt within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Status Updates**: We will provide regular updates on our progress
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days

### Disclosure Policy

- We will work with you to understand and resolve the issue
- We will not take legal action against researchers who report vulnerabilities responsibly
- We will credit you in our security advisory (unless you prefer to remain anonymous)
- We will coordinate public disclosure after the vulnerability is fixed

## Security Considerations

### Vault Token Security

- vssh stores Vault tokens in `~/.vault-token`
- Ensure this file has appropriate permissions (600)
- Tokens are cached for performance but respect Vault TTL
- Use short-lived tokens when possible

### SSH Certificate Security

- SSH certificates are stored in `~/.ssh/vault_signed_*.pub`
- Certificates have limited validity periods
- vssh automatically renews expired certificates
- Private keys remain on the local system

### Configuration Security

- Configuration files may contain sensitive information
- Use environment variables for sensitive values when possible
- Ensure configuration files have appropriate permissions
- Avoid committing configuration files with secrets

### Network Security

- vssh communicates with Vault over HTTPS
- Certificate validation is enforced
- Use Vault's built-in security features (policies, auth methods)

## Best Practices

### For Users

1. **Keep vssh Updated**: Always use the latest version
2. **Secure Configuration**: Protect configuration files and tokens
3. **Monitor Access**: Review Vault audit logs regularly
4. **Principle of Least Privilege**: Use minimal required permissions
5. **Network Security**: Use secure networks and VPNs when possible

### For Administrators

1. **Vault Security**: Follow Vault security best practices
2. **Access Control**: Implement proper RBAC policies
3. **Audit Logging**: Enable and monitor audit logs
4. **Certificate Policies**: Configure appropriate certificate TTLs
5. **Network Segmentation**: Isolate Vault infrastructure

## Security Features

### Authentication

- Multiple authentication methods (Token, UserPass, LDAP, OIDC)
- Secure token storage and management
- Automatic token renewal and validation

### Certificate Management

- Automatic certificate signing and renewal
- Certificate validation and expiration checking
- Secure certificate storage

### Logging

- Structured logging with configurable levels
- No sensitive data in logs
- Debug mode for troubleshooting

## Known Security Considerations

### Token Storage

- Vault tokens are stored in plaintext in `~/.vault-token`
- This follows Vault CLI conventions
- Ensure appropriate file system permissions

### Certificate Caching

- SSH certificates are cached for performance
- Certificates are validated before use
- Expired certificates are automatically renewed

### Configuration Files

- Configuration files may contain sensitive information
- Use environment variables for secrets when possible
- Ensure proper file permissions

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 0.1.6)
- Documented in CHANGELOG.md
- Announced through GitHub releases
- Tagged with security labels

## Contact

For security-related questions or concerns:

1. Create a GitHub issue (for general security questions)
2. Use GitHub's private vulnerability reporting (for vulnerabilities)
3. Contact maintainers directly (for sensitive matters)

## Acknowledgments

We appreciate the security research community and will acknowledge researchers who help improve vssh's security.

---

**Remember**: When in doubt, report it. We'd rather investigate a false positive than miss a real vulnerability.
