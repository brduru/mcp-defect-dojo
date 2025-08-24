# Contributing to MCP DefectDojo

🎉 **Thank you for your interest in contributing!** We welcome contributions from the community and are excited to work with you.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/mcp-defect-dojo.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes: `make test`
6. **Commit** using conventional commits: `git commit -m "feat: add amazing feature"`
7. **Push** to your branch: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

## 🎯 Types of Contributions

### 🐛 Bug Fixes
- Fix existing issues
- Improve error handling
- Performance optimizations

### ✨ New Features
- New MCP tools
- Enhanced DefectDojo integration
- Additional AI agent support

### 📚 Documentation
- README improvements
- Code comments
- Usage examples

### 🧪 Testing
- Unit tests
- Integration tests
- Example validations

## 📋 Development Setup

### Prerequisites
- Go 1.25+
- DefectDojo instance (for testing)
- Make (optional but recommended)

### Local Development
```bash
# Clone the repository
git clone https://github.com/brduru/mcp-defect-dojo.git
cd mcp-defect-dojo

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run with your DefectDojo instance
export DEFECTDOJO_URL="https://your-defectdojo.com"
export DEFECTDOJO_API_KEY="your-api-key"
./bin/mcp-server
```

## 🎨 Code Style

### Go Code Guidelines
- Follow standard Go formatting: `gofmt -s`
- Use meaningful variable names
- Add comments for exported functions
- Follow Go best practices

### Commit Messages
We use [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Features (minor version bump)
feat: add new MCP tool for finding comments
feat(client): implement retry mechanism

# Bug fixes (patch version bump)  
fix: resolve authentication timeout issue
fix(server): handle empty responses correctly

# Documentation (patch version bump)
docs: update installation instructions
docs(readme): add troubleshooting section

# Chores (no version bump)
chore: update dependencies
ci: improve GitHub Actions workflow
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with race detector
make test-race

# Run specific package tests
go test ./internal/defectdojo/...
```

### Writing Tests
- Add tests for new features
- Maintain or improve test coverage
- Include both unit and integration tests
- Mock external dependencies when appropriate

### Test Structure
```go
func TestNewFeature(t *testing.T) {
    // Arrange
    setup := createTestSetup()
    
    // Act
    result, err := NewFeature(setup.input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

## 🔧 Adding New MCP Tools

### 1. Define the Tool
Add your tool definition in `internal/server/mcp.go`:

```go
func (s *MCPServer) registerYourNewTool() {
    tool := &mcp.Tool{
        Name:        "your_new_tool",
        Description: "Description of what your tool does",
        InputSchema: mcp.ToolInputSchema{
            Type: "object",
            Properties: map[string]mcp.ToolInputProperty{
                "param1": {
                    Type:        "string",
                    Description: "Description of param1",
                },
            },
            Required: []string{"param1"},
        },
    }
    
    s.addTool(tool, s.handleYourNewTool)
}
```

### 2. Implement the Handler
```go
func (s *MCPServer) handleYourNewTool(ctx context.Context, params map[string]any) (*mcp.CallToolResult, error) {
    // Extract parameters
    param1, ok := params["param1"].(string)
    if !ok {
        return mcp.NewToolResultError("param1 is required and must be a string"), nil
    }
    
    // Call DefectDojo client
    result, err := s.client.YourNewMethod(ctx, param1)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
    }
    
    // Format response
    return mcp.NewToolResultText(formatResult(result)), nil
}
```

### 3. Add Client Method
Add the DefectDojo API method in `internal/defectdojo/client.go`:

```go
func (c *Client) YourNewMethod(ctx context.Context, param string) (*YourResponse, error) {
    // Implement DefectDojo API call
    url := fmt.Sprintf("%s/api/v2/your-endpoint/?param=%s", c.baseURL, param)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Add authentication, make request, parse response
    // ... implementation details
}
```

### 4. Add Types
Define response types in `pkg/types/defectdojo.go`:

```go
type YourResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    // ... other fields
}
```

### 5. Register the Tool
Add the registration call in `registerTools()`:

```go
func (s *MCPServer) registerTools() {
    s.registerGetFindingsTool()
    s.registerGetFindingDetailTool()
    s.registerMarkFalsePositiveTool()
    s.registerHealthCheckTool()
    s.registerYourNewTool() // Add this line
}
```

## 📝 Pull Request Guidelines

### Before Submitting
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated if needed
- [ ] No merge conflicts with main branch

### Pull Request Template
```markdown
## 🎯 Description
Brief description of changes

## 🔄 Type of Change
- [ ] 🐛 Bug fix
- [ ] ✨ New feature
- [ ] 📚 Documentation update
- [ ] 🧪 Tests
- [ ] 🔧 Maintenance

## 🧪 Testing
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] Manual testing completed

## 📋 Checklist
- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

## 🤝 Code Review Process

1. **Automated Checks**: CI must pass
2. **Maintainer Review**: Core team reviews changes
3. **Feedback**: Address review comments
4. **Approval**: Approved changes are merged
5. **Release**: Changes included in next release

## 🏷️ Release Process

- Contributions are released automatically via conventional commits
- Use correct commit message format for proper versioning
- Breaking changes require `BREAKING CHANGE:` in commit body

## 🆘 Getting Help

- 💬 [GitHub Discussions](https://github.com/brduru/mcp-defect-dojo/discussions)
- 🐛 [Issues](https://github.com/brduru/mcp-defect-dojo/issues)
- 📧 Email maintainers for sensitive topics

## 📜 Code of Conduct

We are committed to providing a welcoming and inclusive environment:

- **Be respectful** of differing viewpoints and experiences
- **Accept constructive criticism** gracefully
- **Focus on what's best** for the community
- **Show empathy** towards other community members

## 🙏 Recognition

Contributors are recognized in:
- Release notes
- README contributors section  
- Git commit history

**Thank you for contributing to MCP DefectDojo!** 🎉
