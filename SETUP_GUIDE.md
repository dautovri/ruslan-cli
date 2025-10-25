# Setting Up ruslan-cli Repository

This guide walks you through setting up the `ruslan-cli` repository from scratch.

## Step 1: Create Repository

```bash
# Create new repository on GitHub
gh repo create dautovri/ruslan-cli --public --description "CLI tool for managing HashiCorp Vault"

# Clone it
cd ~/Developer/Github
git clone https://github.com/dautovri/ruslan-cli.git
cd ruslan-cli
```

## Step 2: Copy Starter Files

```bash
# Copy all template files from gcp-hashicorp-vault repo
cp /path/to/gcp-hashicorp-vault/ruslan-cli-starter/*.template .

# Rename templates (remove .template extension)
for file in *.template; do
    mv "$file" "${file%.template}"
done
```

## Step 3: Set Up Project Structure

```bash
# Create directory structure
mkdir -p cmd
mkdir -p pkg/config
mkdir -p pkg/vault
mkdir -p pkg/discovery
mkdir -p pkg/auth
mkdir -p completions
mkdir -p dist

# Move command files
mv cmd_*.go cmd/

# Move package files
mv pkg_config.go pkg/config/config.go

# Move GitHub workflow
mkdir -p .github/workflows
mv github-workflow-release.yml .github/workflows/release.yml
```

## Step 4: Initialize Go Module

```bash
go mod init github.com/dautovri/ruslan-cli
go mod tidy
```

## Step 5: Create Additional Package Files

### pkg/vault/client.go

```go
package vault

import (
	"fmt"
	
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/dautovri/ruslan-cli/pkg/config"
)

type Client struct {
	*vaultapi.Client
	Config *config.Config
}

func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	
	env := cfg.Environments[cfg.CurrentEnvironment]
	if env == nil {
		return nil, fmt.Errorf("environment not found")
	}
	
	// Discover Vault address if not set
	if env.VaultAddr == "" {
		// TODO: Implement discovery
		return nil, fmt.Errorf("vault address not configured")
	}
	
	vaultCfg := vaultapi.DefaultConfig()
	vaultCfg.Address = env.VaultAddr
	
	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, err
	}
	
	// Load saved token
	// TODO: Implement token loading
	
	return &Client{
		Client: client,
		Config: cfg,
	}, nil
}

// Implement methods: LoginWithToken, GetSecret, PutSecret, etc.
```

### pkg/discovery/discovery.go

```go
package discovery

import (
	"fmt"
	"os/exec"
	"strings"
)

func DiscoverVaultAddress(clusterName, region, namespace, serviceName string) (string, error) {
	// Get kubectl context
	cmd := exec.Command("kubectl", "get", "svc", serviceName,
		"-n", namespace,
		"-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get service IP: %w", err)
	}
	
	ip := strings.TrimSpace(string(output))
	if ip == "" {
		return "", fmt.Errorf("no external IP found for service")
	}
	
	return fmt.Sprintf("http://%s:8200", ip), nil
}
```

## Step 6: Create README.md

```bash
cat > README.md << 'EOF'
# ruslan-cli

CLI tool for managing HashiCorp Vault across multiple environments.

## Installation

### Homebrew
```bash
brew tap dautovri/tap
brew install ruslan-cli
```

### From Source
```bash
go install github.com/dautovri/ruslan-cli@latest
```

## Quick Start

```bash
ruslan-cli env use dev
ruslan-cli login
ruslan-cli secrets list secret/
```

See [documentation](https://github.com/dautovri/gcp-hashicorp-vault/blob/main/docs/RUSLAN_CLI_SPEC.md) for full details.

## License

MIT
EOF
```

## Step 7: Create LICENSE

```bash
cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2025 Ruslan Dautov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF
```

## Step 8: Add Missing Imports to cmd Files

Edit `cmd/root.go` and add missing import:
```go
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
```

## Step 9: Build and Test

```bash
# Install dependencies
go mod download

# Build
make build

# Test
./ruslan-cli --help
./ruslan-cli env list
```

## Step 10: Commit and Push

```bash
git add .
git commit -m "Initial commit: ruslan-cli implementation"
git push origin main
```

## Step 11: Create First Release

```bash
# Tag the release
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0

# This will trigger GitHub Actions to:
# - Build binaries for multiple platforms
# - Create GitHub release
# - Publish to Homebrew tap (if configured)
```

## Step 12: Set Up Homebrew Tap (Optional)

```bash
# Create homebrew-tap repository
gh repo create dautovri/homebrew-tap --public --description "Homebrew tap for dautovri tools"

# Add GitHub secret for tap updates
# Go to ruslan-cli repository settings -> Secrets
# Add HOMEBREW_TAP_GITHUB_TOKEN with access to homebrew-tap repo
```

## Next Steps

1. Implement remaining Vault client methods
2. Add tests
3. Implement auto-discovery logic
4. Add more authentication methods
5. Enhance error handling
6. Add examples directory

## Development

```bash
# Run locally
make run ARGS="env list"

# Run tests
make test

# Build for all platforms
make build-all

# Generate completions
make completions
```

## Troubleshooting

### Build Errors

```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

### Import Errors

Make sure all cmd files have correct package declaration:
```go
package cmd
```

And all use the correct import path:
```go
import "github.com/dautovri/ruslan-cli/pkg/config"
```
