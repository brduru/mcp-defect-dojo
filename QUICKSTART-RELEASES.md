# Quick Start: Releases & Deployment

## 🚀 For Developers

### Making a Release

1. **Automatic Release (Recommended):**
   ```bash
   # For bug fixes (patch ve**Current release pipeline builds:**
- ✅ 5 platform binaries
- ✅ SHA256 checksums
- ✅ Automated changelog
- ✅ Go module compatibility
- ✅ Multi-Go version testing (1.23, 1.24, 1.25)
   git commit -m "fix: resolve authentication timeout"
   git push origin main
   
   # For new features (minor version)
   git commit -m "feat: add finding comments support"
   git push origin main
   
   # For breaking changes (major version)
   git commit -m "feat: redesign API interface
   
   BREAKING CHANGE: endpoint structure changed"
   git push origin main
   ```

2. **Manual Release:**
   - Go to [GitHub Actions](https://github.com/brduru/mcp-defect-dojo/actions)
   - Select "Release" workflow
   - Click "Run workflow"
   - Enter version: `v1.2.3`
   - Click "Run workflow"

### Local Development

```bash
# Build and test
make build
make test
make release  # Build all platforms

# Version info
make version
./bin/mcp-server --version
```

## 📦 For Users

### Go Module Installation

```bash
# Latest version
go get github.com/brduru/mcp-defect-dojo@latest

# Specific version
go get github.com/brduru/mcp-defect-dojo@v1.2.3

# Use in code
import "github.com/brduru/mcp-defect-dojo/pkg/mcpserver"
```

### Binary Installation

```bash
# Linux (amd64)
wget https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-amd64
chmod +x mcp-defect-dojo-linux-amd64
./mcp-defect-dojo-linux-amd64 --version

# macOS (Apple Silicon)
wget https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-darwin-arm64
chmod +x mcp-defect-dojo-darwin-arm64

# Windows
# Download mcp-defect-dojo-windows-amd64.exe from releases page
```

### Verify Downloads

```bash
# Download checksums
wget https://github.com/brduru/mcp-defect-dojo/releases/latest/download/checksums.txt

# Verify integrity
sha256sum -c checksums.txt
```

## 🔄 Release Pipeline

The project uses 3 GitHub Actions workflows:

1. **CI (`ci.yml`)** - Runs on every push/PR
   - Tests across Go 1.21, 1.22, 1.23
   - Linting and formatting checks
   - Build validation

2. **Auto Version (`auto-version.yml`)** - Runs on main branch
   - Analyzes conventional commits
   - Creates semantic version tags
   - Updates CHANGELOG.md

3. **Release (`release.yml`)** - Runs on version tags
   - Builds multi-platform binaries
   - Creates GitHub releases
   - Uploads assets with checksums

## 📋 Commit Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

| Type | Release | Example |
|------|---------|---------|
| `fix:` | Patch | `fix: resolve timeout issue` |
| `feat:` | Minor | `feat: add new finding filter` |
| `feat!:` | Major | `feat!: redesign API` |
| `docs:` | Patch | `docs: update README` |
| `chore:` | None | `chore: update deps` |

## 🎯 Best Practices

### For Contributors
- Use conventional commits for automatic versioning
- Run `make test` before pushing
- Update documentation for new features
- Add examples for new functionality

### For Maintainers
- Review all PRs before merging
- Use semantic versioning consistently
- Keep CHANGELOG.md updated (automatic)
- Test releases before announcement

### For Users
- Pin to specific versions in production
- Verify binary checksums
- Read release notes for breaking changes
- Update dependencies regularly

## 🔧 Troubleshooting

### Common Issues

**Release not created:**
- Check commit message format
- Ensure push is to `main` branch
- Look at GitHub Actions logs

**Build failures:**
- Verify Go version compatibility
- Check all tests pass locally
- Review linting errors

**Version conflicts:**
- Use manual release workflow
- Check for existing tags

### Getting Help

1. Check [existing issues](https://github.com/brduru/mcp-defect-dojo/issues)
2. Review [release documentation](RELEASE.md)
3. Open new issue with `release` label

## 📊 Release Stats

Current release pipeline builds:
- ✅ 5 platform binaries
- ✅ SHA256 checksums
- ✅ Automated changelog
- ✅ Go module compatibility
- ✅ Multi-Go version testing

**Build Targets:**
- `linux/amd64` - Most common Linux servers
- `linux/arm64` - ARM-based servers (e.g., AWS Graviton)
- `darwin/amd64` - Intel Macs
- `darwin/arm64` - Apple Silicon Macs
- `windows/amd64` - Windows systems
