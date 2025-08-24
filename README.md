# MCP DefectDojo

> 🔗 **Connect AI agents to DefectDojo vulnerability management**

[![CI](https://github.com/brduru/mcp-defect-dojo/workflows/CI/badge.svg)](https://github.com/brduru/mcp-defect-dojo/action)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/brduru/mcp-defect-dojo)](https://goreportcard.com/report/github.com/brduru/mcp-defect-dojo)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/brduru/mcp-defect-dojo.svg)](https://pkg.go.dev/github.com/brduru/mcp-defect-dojo)

A [Model Context Protocol](https://spec.modelcontextprotocol.io/) server that enables AI agents to interact with DefectDojo vulnerability management platforms through natural language.

**Compatible with**: Claude Desktop, VS Code Copilot, custom AI agents, and any MCP-compatible tools

## � Quick Start

### For AI Agents (Recommended)

1. **Download the binary**:
   ```bash
   # Linux/macOS
   curl -L https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-amd64 -o mcp-defect-dojo
   chmod +x mcp-defect-dojo
   ```

2. **Configure your AI client**:

   **Claude Desktop** (`~/.claude/claude_desktop_config.json`):
   ```json
   {
     "mcpServers": {
       "defectdojo": {
         "command": "/path/to/mcp-defect-dojo",
         "env": {
           "DEFECTDOJO_URL": "https://your-defectdojo.com",
           "DEFECTDOJO_API_KEY": "your-api-key"
         }
       }
     }
   }
   ```

3. **Start chatting**:
   ```
   You: "Check if DefectDojo is working"
   Claude: ✅ DefectDojo Health Check: HEALTHY
   
   You: "Show me all critical vulnerabilities"
   Claude: Found 5 critical findings...
   ```

### For Go Applications

```bash
go get github.com/brduru/mcp-defect-dojo/pkg/mcpserver
```

```go
package main

import "github.com/brduru/mcp-defect-dojo/pkg/mcpserver"

func main() {
    // Quick setup with API key
    server, err := mcpserver.NewServerWithAPIKey("your-api-key")
    if err != nil {
        panic(err)
    }
    
    // Run the server
    if err := server.Run(context.Background()); err != nil {
        panic(err)
    }
}
```

## 🛠️ Available Tools

| Tool | Description | Example |
|------|-------------|---------|
| `defectdojo_health_check` | Verify connectivity | *"Is DefectDojo online?"* |
| `get_defectdojo_findings` | Search vulnerabilities | *"Show me all critical findings"* |
| `get_finding_detail` | Get finding details | *"Get details for finding #123"* |
| `mark_finding_false_positive` | Mark false positives | *"Mark finding #456 as false positive"* |

### Example Conversations

```
🧑: "Check if DefectDojo is working"
🤖: ✅ DefectDojo Health Check: HEALTHY
   Connection successful to https://your-defectdojo.com
   API v2 is responsive and accessible.

🧑: "Show me the most critical vulnerabilities"  
🤖: Found 3 critical findings:
   
   1. [Critical] SQL Injection in Authentication (ID: 456)
      Status: Active, Verified: true
      
   2. [Critical] Remote Code Execution via Upload (ID: 789)  
      Status: Active, Verified: false
```
## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DEFECTDOJO_URL` | DefectDojo base URL | `http://localhost:8080` | ✅ |
| `DEFECTDOJO_API_KEY` | API authentication key | - | ✅ |
| `DEFECTDOJO_API_VERSION` | API version | `v2` | ❌ |

### Configuration Methods

```go
// Method 1: Environment variables (recommended for AI agents)
server, err := mcpserver.NewServer()

// Method 2: Direct API key
server, err := mcpserver.NewServerWithAPIKey("your-api-key")

// Method 3: Full configuration
server, err := mcpserver.NewServerWithSettings(mcpserver.DefectDojoSettings{
    BaseURL:    "https://defectdojo.company.com",
    APIKey:     "your-api-key",
    APIVersion: "v2",
})
```
## 📦 Installation

### Pre-built Binaries

| Platform | Download |
|----------|----------|
| Linux (x64) | [mcp-defect-dojo-linux-amd64](https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-amd64) |
| Linux (ARM64) | [mcp-defect-dojo-linux-arm64](https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-arm64) |
| macOS (Intel) | [mcp-defect-dojo-darwin-amd64](https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-darwin-amd64) |
| macOS (Apple Silicon) | [mcp-defect-dojo-darwin-arm64](https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-darwin-arm64) |
| Windows (x64) | [mcp-defect-dojo-windows-amd64.exe](https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-windows-amd64.exe) |

### From Source

```bash
git clone https://github.com/brduru/mcp-defect-dojo.git
cd mcp-defect-dojo
make build
```

### Go Module

```bash
go get github.com/brduru/mcp-defect-dojo/pkg/mcpserver@latest
```

## 🔧 Development

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./pkg/mcpserver -v
```

**Current Test Coverage:**
- `pkg/mcpserver`: 32.9%
- `pkg/types`: 100%
- `internal/config`: 80%
- `internal/defectdojo`: 86.9%

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run locally
make run
```

## 📚 Documentation

- **[GoDoc API Reference](https://pkg.go.dev/github.com/brduru/mcp-defect-dojo)** - Complete API documentation
- **[Examples](examples/)** - Usage examples and integration patterns
- **[Model Context Protocol](https://spec.modelcontextprotocol.io/)** - Learn about MCP
- **[DefectDojo API](https://demo.defectdojo.org/api/v2/)** - DefectDojo API documentation

## 🤝 Contributing

We welcome contributions! Please see our [contributing guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Connect your AI agents to DefectDojo and automate vulnerability management** 🚀
