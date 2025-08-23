# Conventional Commits Examples

This file provides examples of conventional commit messages for automatic versioning.

## Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Examples by Release Type

### Patch Releases (1.0.0 → 1.0.1)

Bug fixes and minor improvements:

```bash
fix: resolve DefectDojo authentication timeout
fix: correct finding status enumeration  
fix(client): handle empty response from API properly
docs: fix typo in README installation section
chore: update dependencies to latest versions
```

### Minor Releases (1.0.0 → 1.1.0)

New features (backward compatible):

```bash
feat: add support for finding comments retrieval
feat: implement bulk finding operations
feat(mcp): add new tool for product listing
perf: optimize DefectDojo API response parsing
```

### Major Releases (1.0.0 → 2.0.0)

Breaking changes:

```bash
feat: redesign MCP tool interface

BREAKING CHANGE: Tool response format has changed. 
All tools now return structured JSON instead of plain text.

feat!: remove deprecated SSE transport support

BREAKING CHANGE: SSE transport has been removed. 
Use stdio or in-process transport instead.
```

## Multi-line Commit Example

```bash
feat: add advanced finding filtering capabilities

- Support filtering by severity, status, and date range
- Add support for custom query parameters  
- Implement pagination for large result sets
- Add validation for filter parameters

Closes #123
Co-authored-by: Jane Developer <jane@example.com>
```

## Special Cases

### Revert Commits
```bash
revert: feat: add experimental SSE support

This reverts commit 1234567890abcdef due to stability issues.
```

### Multiple Types
```bash
feat: add health check endpoint
fix: resolve timeout in existing endpoints
docs: update API documentation

BREAKING CHANGE: Health check endpoint moved to /health
```

## Scope Examples

Use scopes to specify the area of change:

```bash
feat(client): add retry mechanism for failed requests
fix(server): resolve memory leak in tool registration  
docs(readme): update installation instructions
test(integration): add DefectDojo API integration tests
ci: update GitHub Actions to Go 1.25
```

## Release Automation

These commit messages will trigger:

- **fix/docs/perf/refactor**: Patch version bump
- **feat**: Minor version bump  
- **BREAKING CHANGE**: Major version bump
- **chore/ci/test**: No version bump

## Tools for Validation

Install commitizen for guided commit messages:

```bash
npm install -g commitizen cz-conventional-changelog
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc

# Use 'git cz' instead of 'git commit'
git cz
```

Or use conventional-commit tools:

```bash
# Install conventional-commit CLI
npm install -g @commitlint/cli @commitlint/config-conventional

# Add commit-msg hook
echo "module.exports = {extends: ['@commitlint/config-conventional']}" > commitlint.config.js
```
