# Installing ruslan-cli via Homebrew

## Prerequisites

1. **Create Homebrew Tap Repository**
   ```bash
   gh repo create dautovri/homebrew-tap --public --description "Homebrew tap for dautovri tools"
   cd ~/Developer/Github
   git clone https://github.com/dautovri/homebrew-tap.git
   ```

2. **Set Up GitHub Token for Tap Updates**
   
   The release workflow needs a token to push formula updates to your tap repository.
   
   a. Create a Personal Access Token:
   - Go to https://github.com/settings/tokens/new
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Expiration: Choose appropriate duration
   - Scopes needed:
     - `repo` (Full control of private repositories)
     - `write:packages` (if using GitHub Packages)
   
   b. Add token to ruslan-cli repository secrets:
   - Go to https://github.com/dautovri/ruslan-cli/settings/secrets/actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: Paste the token you created
   - Click "Add secret"

## Creating Your First Release

1. **Commit and push all changes**:
   ```bash
   cd /Users/rd/Developer/Github/ruslan-cli
   git add .
   git commit -m "feat: initial release with Homebrew support"
   git push origin main
   ```

2. **Create and push a version tag**:
   ```bash
   git tag -a v0.1.0 -m "Initial release: ruslan-cli v0.1.0

   Features:
   - Multi-environment support (dev/prod)
   - Vault authentication (token, userpass, approle)
   - KV v2 secrets management (list, get, put, delete)
   - JSON, YAML, and table output formats
   - Shell completions (bash, zsh, fish)
   "
   git push origin v0.1.0
   ```

3. **Monitor the release**:
   - GitHub Actions will automatically:
     - Build binaries for Linux, macOS, Windows (amd64/arm64)
     - Create a GitHub Release with assets
     - Generate and push Homebrew Formula to your tap
   
   Check progress: https://github.com/dautovri/ruslan-cli/actions

## Installing via Homebrew

Once the first release is complete, users can install:

```bash
# Add the tap
brew tap dautovri/tap

# Install ruslan-cli
brew install ruslan-cli

# Verify installation
ruslan-cli --help
ruslan-cli env list
```

## Updating the Formula

Future releases are automatic! Just create and push a new tag:

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

GoReleaser will:
- Build new binaries
- Create GitHub release
- Update Homebrew formula automatically

## Manual Testing Before Release

Test the build locally:

```bash
# Install GoReleaser
brew install goreleaser

# Test the release process (without publishing)
goreleaser release --snapshot --clean

# Check the generated artifacts
ls -la dist/
```

## Troubleshooting

### Token Issues
If the release fails with authentication errors:
1. Verify `HOMEBREW_TAP_GITHUB_TOKEN` secret exists
2. Check token has correct permissions
3. Ensure token hasn't expired

### Formula Not Updated
If Homebrew formula isn't created:
1. Verify `homebrew-tap` repository exists
2. Check GitHub Actions logs for errors
3. Ensure token has write access to tap repo

### Build Failures
If builds fail:
1. Check `go.mod` is valid: `go mod tidy`
2. Verify tests pass: `go test ./...`
3. Check GoReleaser config: `goreleaser check`

## Current Status

✅ `.goreleaser.yml` configured
✅ Release workflow ready (`.github/workflows/release.yml`)
✅ Shell completions support added
✅ Multi-platform builds (Linux/macOS/Windows, amd64/arm64)

🔄 **Next Steps:**
1. Create `homebrew-tap` repository
2. Add `HOMEBREW_TAP_GITHUB_TOKEN` secret
3. Create first release tag (`v0.1.0`)
4. Test installation: `brew tap dautovri/tap && brew install ruslan-cli`

## Alternative Installation (Before Homebrew)

Until Homebrew is set up, users can install from source:

```bash
# Clone repository
git clone https://github.com/dautovri/ruslan-cli.git
cd ruslan-cli

# Build and install
make install

# Or just build locally
make build
sudo cp ruslan-cli /usr/local/bin/
```

Or download from GitHub Releases (after first tag):
```bash
# Download latest release
VERSION="v0.1.0"
OS="darwin"  # or "linux" or "windows"
ARCH="arm64"  # or "amd64"

curl -L "https://github.com/dautovri/ruslan-cli/releases/download/${VERSION}/ruslan-cli_${VERSION}_${OS}_${ARCH}.tar.gz" | tar xz

# Install
sudo mv ruslan-cli /usr/local/bin/
ruslan-cli --help
```
