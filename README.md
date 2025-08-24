# MCP DefectDojo

🔗 **Connect AI agents to DefectDojo vulnerability management platform**

[![CI](https://github.com/brduru/mcp-defect-dojo/workflows/CI/badge.svg)](https://github.com/brduru/mcp-defect-dojo/action)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/brduru/mcp-defect-dojo)](https://goreportcard.com/report/github.com/brduru/mcp-defect-dojo)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/brduru/mcp-defect-dojo)](https://github.com/brduru/mcp-defect-dojo/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/brduru/mcp-defect-dojo/total.svg)](https://github.com/brduru/mcp-defect-dojo/releases)
[![GitHub stars](https://img.shields.io/github/stars/brduru/mcp-defect-dojo?style=social)](https://github.com/brduru/mcp-defect-dojo/stargazers)


**MCP DefectDojo** is a Model Context Protocol (MCP) integration that enables AI agents to interact with DefectDojo vulnerability management platforms. Use it to automate security workflows, analyze vulnerabilities, and manage findings through natural language AI interactions.

> 🤖 **Perfect for**: Claude Desktop, VS Code Copilot, custom AI agents, and any MCP-compatible tools

## ✨ What You Can Do

**🔍 Query Vulnerabilities:**
- "Show me all critical findings from the last week"
- "Get details about finding #123"
- "List all unverified vulnerabilities in my project"

**🎯 Manage Findings:**
- "Mark finding #456 as false positive - it's test data"
- "Show me all findings for the authentication module"
- "Check if DefectDojo is accessible"

**� Automate Workflows:**
- Filter findings by severity, status, or test
- Bulk analyze vulnerability trends
- Integrate with CI/CD pipelines

## 🎯 Available Tools

| Tool | Description | Use Case |
|------|-------------|----------|
| `defectdojo_health_check` | Verify connectivity | "Is DefectDojo online?" |
| `get_defectdojo_findings` | Search vulnerabilities | "Show critical findings" |
| `get_finding_detail` | Get full finding info | "Details for finding #123" |
| `mark_finding_false_positive` | Mark as false positive | "This is a false alarm" |

## � Quick Start

### 1. For AI Agents (Claude Desktop, VS Code, etc.)

**Download & Run:**
```bash
# Download latest release
curl -L https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-amd64 -o mcp-defect-dojo
chmod +x mcp-defect-dojo

# Set your DefectDojo connection
export DEFECTDOJO_URL="https://your-defectdojo.com"
export DEFECTDOJO_API_KEY="your-api-key"

# Start MCP server
./mcp-defect-dojo
```

**Claude Desktop Configuration:**
Add to your `claude_desktop_config.json`:
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

### 2. For Go Applications

**Install Library:**
```bash
go get github.com/brduru/mcp-defect-dojo/pkg/mcpserver
```

**Simple Integration:**
```go
import "github.com/brduru/mcp-defect-dojo/pkg/mcpserver"

// Quick setup with API key
server, err := mcpserver.NewServerWithAPIKey("your-api-key")

// Or full configuration
server, err := mcpserver.NewServerWithSettings(mcpserver.DefectDojoSettings{
    BaseURL: "https://defectdojo.company.com",
    APIKey:  "your-api-key",
})

// Use with in-process client
client, err := client.NewInProcessClient(server.GetMCPServer())
```

## 💬 Example Conversations

**With Claude Desktop:**
```
You: "Check if DefectDojo is working"
Claude: I'll check the DefectDojo health status for you.
✅ DefectDojo Health Check: HEALTHY
Connection successful to https://your-defectdojo.com
API v2 is responsive and accessible.

You: "Show me the most critical vulnerabilities"
Claude: I'll retrieve the critical findings from DefectDojo.

Found 15 findings (showing 10):

1. [Critical] SQL Injection in Authentication Module (ID: 456)
   Active: true, Verified: true, False Positive: false
   Description: SQL injection vulnerability in login endpoint

2. [Critical] Remote Code Execution via File Upload (ID: 789)
   Active: true, Verified: false, False Positive: false
   Description: Unrestricted file upload allows RCE
```

**With MCP-compatible Tools:**
```bash
# Using any MCP client via stdio
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "defectdojo_health_check", "arguments": {}}}' | ./mcp-defect-dojo

# List available tools
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./mcp-defect-dojo
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DEFECTDOJO_URL` | DefectDojo base URL | `http://localhost:8080` | ✅ |
| `DEFECTDOJO_API_KEY` | API authentication key | - | ✅ |
| `DEFECTDOJO_API_VERSION` | API version | `v2` | ❌ |
| `LOG_LEVEL` | Logging level | `info` | ❌ |

### AI Agent Setup

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

**VS Code MCP Extension**:
```json
{
  "mcp.servers": [
    {
      "name": "defectdojo",
      "command": "/path/to/mcp-defect-dojo",
      "env": {
        "DEFECTDOJO_URL": "https://your-defectdojo.com",
        "DEFECTDOJO_API_KEY": "your-api-key"
      }
    }
  ]
}
```

## 🛠️ Tool Reference

### `defectdojo_health_check`
**Purpose:** Verify DefectDojo connectivity and status  
**Parameters:** None  
**Example:** "Is DefectDojo online?"

### `get_defectdojo_findings`
**Purpose:** Search and filter vulnerability findings  
**Parameters:**
- `limit` (number): Max results to return (default: 10)
- `offset` (number): Pagination offset (default: 0)
- `active_only` (boolean): Show only active findings (default: true)
- `severity` (string): Filter by severity (Critical, High, Medium, Low, Info)
- `test` (number): Filter by specific test ID

**Example:** "Show me all critical findings"

### `get_finding_detail`
**Purpose:** Get comprehensive information about a specific finding  
**Parameters:**
- `finding_id` (number): The ID of the finding to retrieve

**Example:** "Get details for finding #123"

### `mark_finding_false_positive`
**Purpose:** Mark a finding as false positive with justification  
**Parameters:**
- `finding_id` (number): The ID of the finding to mark
- `justification` (string): Reason for marking as false positive
- `notes` (string, optional): Additional notes

**Example:** "Mark finding #456 as false positive - it's test data"

## � Installation Options

### Pre-built Binaries
Download ready-to-use binaries for your platform:

| Platform | Download Link |
|----------|---------------|
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

## 🤝 Support & Community

### Documentation
- [Model Context Protocol](https://spec.modelcontextprotocol.io/) - Learn about MCP
- [DefectDojo API](https://demo.defectdojo.org/api/v2/) - DefectDojo API documentation
- [Examples](examples/) - Complete usage examples

### Issues & Questions
- 🐛 [Report bugs](https://github.com/brduru/mcp-defect-dojo/issues/new?template=bug_report.md)
- 💡 [Request features](https://github.com/brduru/mcp-defect-dojo/issues/new?template=feature_request.md)
- ❓ [Ask questions](https://github.com/brduru/mcp-defect-dojo/discussions)

### Contributing
We welcome contributions! See our [contributing guide](CONTRIBUTING.md) for details.

## � License

MIT License - see [LICENSE](LICENSE) file for details.

---
*Connect your AI agents to DefectDojo and automate vulnerability management*
