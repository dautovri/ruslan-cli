# ruslan-cli

[![CI](https://github.com/dautovri/ruslan-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/dautovri/ruslan-cli/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dautovri/ruslan-cli)](https://goreportcard.com/report/github.com/dautovri/ruslan-cli)
[![codecov](https://codecov.io/gh/dautovri/ruslan-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/dautovri/ruslan-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

CLI tool for managing HashiCorp Vault across multiple environments (dev/prod).

## Features

- 🌍 Multi-environment support (dev/prod)
- 🔐 Multiple authentication methods (token, userpass, approle)
- 🔑 Secret management (list, get, put, delete)
- 📊 Multiple output formats (table, JSON, YAML)
- 🍺 Easy installation via Homebrew

## Installation

### Homebrew (Recommended)
```bash
# Add tap
brew tap dautovri/tap

# Install
brew install ruslan-cli

# Verify
ruslan-cli --help
```

### From Source
```bash
go install github.com/dautovri/ruslan-cli@latest
```

### Download Binary
Download from [GitHub Releases](https://github.com/dautovri/ruslan-cli/releases):
```bash
# macOS (Apple Silicon)
curl -L "https://github.com/dautovri/ruslan-cli/releases/latest/download/ruslan-cli_darwin_arm64.tar.gz" | tar xz
sudo mv ruslan-cli /usr/local/bin/

# macOS (Intel)
curl -L "https://github.com/dautovri/ruslan-cli/releases/latest/download/ruslan-cli_darwin_amd64.tar.gz" | tar xz
sudo mv ruslan-cli /usr/local/bin/

# Linux
curl -L "https://github.com/dautovri/ruslan-cli/releases/latest/download/ruslan-cli_linux_amd64.tar.gz" | tar xz
sudo mv ruslan-cli /usr/local/bin/
```

### Build Locally
```bash
git clone https://github.com/dautovri/ruslan-cli.git
cd ruslan-cli
make build
sudo make install
```

## Quick Start

```bash
# List available environments
ruslan-cli env list

# Switch to dev environment
ruslan-cli env use dev

# Login with token
ruslan-cli login --method=token

# List secrets
ruslan-cli secrets list secret/

# Get a secret
ruslan-cli secrets get secret/myapp/config

# Put a secret
ruslan-cli secrets put secret/myapp/config key=value

# Delete a secret
ruslan-cli secrets delete secret/myapp/config
```

## Configuration

ruslan-cli reads configuration from `~/.ruslan-cli/config.yaml` and `.ruslan-cli.yaml` in your infrastructure repository.

Example `.ruslan-cli.yaml`:
```yaml
environments:
  dev:
    project_id: "homework-475918"
    cluster_name: "dev-gke-cluster"
    region: "us-central1"
    namespace: "vault"
    service_name: "vault"
    vault_addr: "https://vault-dev.dautov.dev"
  prod:
    project_id: "homework-475918"
    cluster_name: "prod-gke-cluster"
    region: "us-central1"
    namespace: "vault"
    service_name: "vault"
    vault_addr: "https://vault.dautov.dev"
    cluster_name: "prod-gke-cluster"
    region: "us-central1"
    namespace: "vault"
    service_name: "vault"
```

## Commands

### Environment Management
- `env list` - List all environments
- `env use <name>` - Switch to an environment
- `env current` - Show current environment
- `env info` - Show environment details

### Authentication
- `login --method=token` - Login with token
- `login --method=userpass` - Login with username/password
- `login --method=approle` - Login with AppRole
- `logout` - Clear saved credentials
- `auth status` - Show authentication status

### Secret Management
- `secrets list <path>` - List secrets at path
- `secrets get <path>` - Read a secret
- `secrets put <path> key=value` - Write a secret
- `secrets delete <path>` - Delete a secret

## Development

### Quick Start

```bash
# Build the CLI
make build

# Run all tests
make test

# Install locally
make install
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests with race detector (recommended)
make test-race

# Run CI-style tests (race + coverage)
make test-ci

# Run benchmarks
make bench

# Show coverage percentage
make coverage
```

### Code Quality

This project uses `golangci-lint` with 25+ linters for code quality enforcement.

```bash
# Run all linters
make lint

# Auto-fix linting issues where possible
make lint-fix

# Verify code quality (format + lint + test)
make verify
```

### CI/CD Pipeline

Every push and pull request automatically triggers:

- ✅ **Linting**: golangci-lint, gofmt, go mod tidy validation
- ✅ **Testing**: Matrix testing (Ubuntu/macOS × Go 1.21/1.22) with race detector
- ✅ **Building**: Cross-platform builds (Linux, macOS, Windows)
- ✅ **Integration**: CLI command validation tests
- ✅ **Security**: gosec static analysis and trivy vulnerability scanning
- ✅ **Coverage**: Automated coverage reporting via Codecov

Run the full CI suite locally before pushing:

```bash
make ci
```

See the [CI workflow](.github/workflows/ci.yml) for complete details.

### Build Artifacts

```bash
# Build for all platforms
make build-all

# Generate shell completions
make completions
```

## License

MIT
