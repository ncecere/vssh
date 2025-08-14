# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Future features will be listed here

## [0.1.5] - 2025-01-13

### Fixed
- Disabled Go module caching for build jobs to eliminate tar extraction conflicts
- Temporarily removed caching to prevent "Cannot open: File exists" errors in parallel builds
- Kept caching enabled for test job since it runs in isolation
- This should eliminate the 60+ tar extraction errors that persisted in v0.1.4

## [0.1.4] - 2025-01-13

### Fixed
- Fixed release workflow cache conflicts causing 60+ build errors
- Used unique cache keys per build matrix combination to prevent conflicts
- Each build job now has isolated Go module cache
- Fixed tar extraction conflicts that were causing "Cannot open: File exists" errors
- Improved build isolation to prevent parallel job interference

## [0.1.3] - 2025-01-13

### Fixed
- Fixed release workflow build failures
- Removed problematic GitHub Packages publishing job
- Added binary cleanup to prevent file conflicts
- Used unique artifact names to avoid naming collisions
- Simplified workflow for better reliability

## [0.1.2] - 2025-01-13

### Fixed
- Fixed duplicate workflow runs when pushing version tags
- CI workflow now only runs on pull requests
- Release workflow handles all tag-based builds and testing

## [0.1.1] - 2025-01-13

### Fixed
- Fixed CI workflow to only run on tags instead of every commit
- Fixed unit tests to match new design where username is used as Vault role

## [0.1.0] - 2025-01-13

### Added
- Initial release of vssh
- Multiple Vault authentication methods (Token, UserPass, LDAP, OIDC)
- Automatic SSH certificate management
- Username-based Vault role mapping
- Multi-user SSH key support
- Certificate caching and validation
- Comprehensive logging and debugging
- Cross-platform support (Linux, macOS, Windows)
- GitHub Actions CI/CD pipeline
- Installation scripts for all platforms
- Comprehensive documentation

### Features
- **Vault Integration**: Seamless integration with HashiCorp Vault SSH secrets engine
- **Authentication Methods**: Support for Token, Username/Password, LDAP, and OIDC authentication
- **Certificate Management**: Automatic certificate signing, validation, and renewal
- **SSH Compatibility**: Works as a drop-in replacement for SSH with additional Vault features
- **Multi-Platform**: Native binaries for Linux (AMD64/ARM64), macOS (Intel/Apple Silicon), and Windows (AMD64/ARM64)
- **Configuration**: Flexible YAML-based configuration with environment variable support
- **Logging**: Structured logging with debug and verbose modes
- **Security**: Secure token storage and certificate lifecycle management

### Documentation
- Complete user documentation in README.md
- Configuration reference in CONFIG.md
- Developer guide in DEVELOPMENT.md
- Installation instructions for all platforms
- Troubleshooting guides and examples

### Infrastructure
- GitHub Actions workflows for CI/CD
- Automated testing across multiple platforms
- Security scanning and vulnerability checks
- Cross-platform binary builds
- Automated release creation
- Installation scripts for Linux, macOS, and Windows

## [1.0.0] - TBD

### Added
- Initial stable release

---

## Release Notes Template

When creating a new release, use this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Bug fixes

### Security
- Security improvements
```

## Version History

- **v1.0.0**: Initial stable release (planned)
- **v0.1.0**: Initial development version

## Migration Guide

### From v0.x to v1.0

No migration required for initial release.

## Breaking Changes

None yet.

## Security Updates

All security updates will be documented here with CVE numbers if applicable.

## Deprecation Notices

No deprecations yet.
