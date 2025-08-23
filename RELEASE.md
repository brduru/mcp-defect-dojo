# Release Process

This document outlines the release process for the MCP DefectDojo server.

## Overview

The project uses **semantic versioning** (semver) and provides automated releases through GitHub Actions with both **automatic** and **manual** workflows.

## Release Types

### 🤖 Automatic Releases (Recommended)

Triggered automatically when commits are pushed to `main` branch using **conventional commits**.

**Commit Message Format:**
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat:` → **Minor** version bump (new features)
- `fix:` → **Patch** version bump (bug fixes)  
- `docs:` → **Patch** version bump (documentation)
- `chore:` → No version bump
- `BREAKING CHANGE:` → **Major** version bump (breaking changes)

**Examples:**
```bash
# Patch release (1.0.0 → 1.0.1)
git commit -m "fix: resolve authentication timeout issue"

# Minor release (1.0.0 → 1.1.0)  
git commit -m "feat: add new finding filtering capabilities"

# Major release (1.0.0 → 2.0.0)
git commit -m "feat: redesign API interface

BREAKING CHANGE: API endpoint structure has changed"
```

### 🚀 Manual Releases

Triggered manually via GitHub Actions UI for specific version releases.

1. Go to **Actions** → **Release** workflow
2. Click **Run workflow**
3. Enter version (e.g., `v1.2.3`)
4. Click **Run workflow**

## What Gets Released

### 📦 Go Module
- Tagged release on GitHub for `go get` compatibility
- Semantic version tag (e.g., `v1.2.3`)

### 📁 Binary Assets
- Multi-platform binaries:
  - `mcp-defect-dojo-linux-amd64`
  - `mcp-defect-dojo-linux-arm64`
  - `mcp-defect-dojo-darwin-amd64` (Intel Mac)
  - `mcp-defect-dojo-darwin-arm64` (Apple Silicon)
  - `mcp-defect-dojo-windows-amd64.exe`
- SHA256 checksums (`checksums.txt`)

### 📋 Release Notes
- Automatically generated changelog
- Installation instructions
- Usage examples

## Local Development

### Building Locally
```bash
# Build for current platform
make build

# Build for all platforms  
make release

# Show version info
make version

# Run tests
make test
```

### Version Information
```bash
# Check current version
./bin/mcp-server --version

# Or using make
make version
```

## CI/CD Pipeline

### 🔄 Continuous Integration (`ci.yml`)
- **Triggers:** Push to `main`, Pull Requests
- **Actions:** Tests, linting, multi-Go version compatibility
- **Matrix:** Go 1.23, 1.24, 1.25

### 📈 Auto Versioning (`auto-version.yml`)  
- **Triggers:** Push to `main` (conventional commits)
- **Actions:** Analyze commits, bump version, create tag
- **Tools:** semantic-release

### 🎁 Release (`release.yml`)
- **Triggers:** Version tags (`v*`), manual workflow
- **Actions:** Build binaries, create GitHub release
- **Artifacts:** Multi-platform binaries + checksums

## Release Checklist

### Before Release
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Examples validated
- [ ] Version updated in code (automatic)

### During Release
- [ ] Use conventional commit messages
- [ ] Or trigger manual release with correct version
- [ ] Verify CI/CD pipeline completion

### After Release
- [ ] Verify GitHub release created
- [ ] Test `go get` with new version
- [ ] Download and test binary assets
- [ ] Update any dependent projects

## Troubleshooting

### Common Issues

**Release not triggered:**
- Check commit message format for conventional commits
- Ensure push is to `main` branch
- Verify no `[skip ci]` in commit message

**Build failures:**
- Check Go version compatibility
- Verify all tests pass locally
- Review action logs for specific errors

**Version conflicts:**
- Use manual release workflow to override
- Check for existing tags with same version

### Manual Recovery

If automatic release fails, you can:

1. **Fix and re-trigger:**
   ```bash
   git commit -m "fix: resolve release issue"
   git push origin main
   ```

2. **Manual release:**
   - Use GitHub Actions UI with specific version
   - Or create tag manually:
     ```bash
     git tag v1.2.3
     git push origin v1.2.3
     ```

## Support

For release-related issues:
1. Check [GitHub Actions](../../actions) for workflow status
2. Review [conventional commits](https://www.conventionalcommits.org/) specification
3. Open issue with `release` label
