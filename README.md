# MCP DefectDojo Server

A modern Go library and MCP (Model Context Protocol) server for DefectDojo integration, built with [mcp-go](https://github.com/mark3labs/mcp-go).

[![CI](https://github.com/brduru/mcp-defect-dojo/workflows/CI/badge.svg)](https://github.com/brduru/mcp-defect-dojo/actions/workflows/ci.yml)
[![Auto Version](https://github.com/brduru/mcp-defect-dojo/workflows/Auto%20Version/badge.svg)](https://github.com/brduru/mcp-defect-dojo/actions/workflows/auto-version.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/brduru/mcp-defect-dojo)](https://github.com/brduru/mcp-defect-dojo/releases/latest)

## 🚀 Overview

This project provides MCP tools for AI agents to interact with DefectDojo vulnerability management platform. It supports multiple transport methods and can be used both as a standalone server and as a Go library.

### ✨ Features

**🛠️ MCP Tools:**
- `defectdojo_health_check` - Verify DefectDojo connectivity and status
- `get_defectdojo_findings` - Retrieve vulnerability findings with advanced filtering
- `get_finding_detail` - Get comprehensive information about a specific finding
- `mark_finding_false_positive` - Mark findings as false positive with justification

**🔌 Transport Support:**
- **📡 Stdio** - For subprocess/sidecar usage (recommended for AI agents)
- **🔧 In-Process** - For direct Go library integration via `go get`

**🏗️ Architecture:**
- Clean, modular Go architecture
- Full DefectDojo API v2 integration
- Programmatic configuration (no environment files needed)
- Comprehensive error handling
- Type-safe MCP responses

## 📦 Installation

### Option 1: As a Go Library (In-Process Transport)

For direct integration into your Go applications:

```bash
go get github.com/brduru/mcp-defect-dojo/pkg/mcpserver
```

**Example Usage:**

```go
package main

import (
    "context"
    "log"
    
    "github.com/brduru/mcp-defect-dojo/pkg/mcpserver"
    "github.com/mark3labs/mcp-go/client"
    "github.com/mark3labs/mcp-go/mcp"
)

func main() {
    // Option 1: Simple usage with API key (uses localhost:8080)
    server, err := mcpserver.NewServerWithAPIKey("your-api-key")
    
    // Option 2: Full control with custom DefectDojo settings
    server, err := mcpserver.NewServerWithSettings(mcpserver.DefectDojoSettings{
        BaseURL:    "https://defectdojo.company.com",
        APIKey:     "your-api-key",
        APIVersion: "v2", // optional
    })
    
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    // Create in-process client
    mcpClient, err := client.NewInProcessClient(server.GetMCPServer())
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Call DefectDojo tools directly
    result, err := mcpClient.CallTool(context.Background(), mcp.CallToolRequest{
        Params: mcp.CallToolParams{
            Name: "defectdojo_health_check",
            Arguments: map[string]any{},
        },
    })
    
    if err != nil {
        log.Fatalf("Tool call failed: %v", err)
    }
    
    log.Printf("Health check result: %v", result.Content)
}
```

> **💡 Usage Options:** 
> - `NewServerWithAPIKey()` - Simple usage with localhost:8080
> - `NewServerWithSettings()` - Full control over DefectDojo connection
```
```

### Option 2: As a Standalone Binary (Stdio Transport)

For use with AI agents and external tools:

**Download pre-built binary:**
```bash
# Download latest release for your platform
wget https://github.com/brduru/mcp-defect-dojo/releases/latest/download/mcp-defect-dojo-linux-amd64
chmod +x mcp-defect-dojo-linux-amd64

# Configure environment
export DEFECTDOJO_URL="https://your-defectdojo.com"
export DEFECTDOJO_API_KEY="your-api-key"

# Run as MCP server
./mcp-defect-dojo-linux-amd64
```

**Or build from source:**
```bash
git clone https://github.com/brduru/mcp-defect-dojo.git
cd mcp-defect-dojo
make build

# Configure environment
export DEFECTDOJO_URL="https://your-defectdojo.com"
export DEFECTDOJO_API_KEY="your-api-key"

# Run as MCP server
./bin/mcp-server

# Check version
./bin/mcp-server --version
```
# Edit .env file with your DefectDojo instance details
vim .env
```
**Use with AI agents:**
```bash
# Example with Claude Desktop or other MCP clients
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./bin/mcp-server
```

## ⚙️ Configuration

The server uses sensible defaults that work out-of-the-box for most DefectDojo setups. Configuration is done via environment variables for DefectDojo connection settings:

### Environment Variables

**DefectDojo Connection (configurable):**
- `DEFECTDOJO_URL` - DefectDojo base URL (default: `http://localhost:8080`)
- `DEFECTDOJO_API_KEY` - DefectDojo API key for authentication
- `DEFECTDOJO_API_VERSION` - DefectDojo API version (default: `v2`)

**Debugging (configurable):**
- `LOG_LEVEL` - Logging level: `debug`, `info`, `warn`, `error` (default: `info`)
- `LOG_FORMAT` - Log format: `text`, `json` (default: `text`)

**Server Identity (fixed - cannot be overridden):**
- **Server Name**: `mcp-defect-dojo-server` 
- **Server Version**: `v0.1.0`
- **Instructions**: DefectDojo MCP integration tools

### Default Configuration

The server comes with these built-in defaults:
- **DefectDojo URL**: `http://localhost:8080` 
- **API Version**: `v2`
- **Request Timeout**: `30 seconds`
- **Server Name**: `mcp-defect-dojo-server` (fixed)
- **Server Version**: `v0.1.0` (fixed)
- **Log Level**: `info`

## 🎯 Quick Start

### 1. Health Check
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "defectdojo_health_check",
    "arguments": {}
  }
}
```

### 2. Get Findings
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_defectdojo_findings",
    "arguments": {
      "limit": 10,
      "active_only": true,
      "severity": "High"
    }
  }
}
```

### 3. Get Finding Details
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_finding_detail",
    "arguments": {
      "finding_id": 123
    }
  }
}
```

### 4. Mark False Positive
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "mark_finding_false_positive",
    "arguments": {
      "finding_id": 123,
      "justification": "Not applicable in this context",
      "notes": "Reviewed by security team"
    }
  }
}
```

## 🏗️ Project Structure

```
mcp-defect-dojo/
├── cmd/
│   └── mcp-server/           # Standalone MCP server binary
├── pkg/
│   ├── mcpserver/            # Public API for go get usage
│   └── types/                # Shared types and data structures
├── internal/
│   ├── server/               # MCP server implementation
│   ├── defectdojo/           # DefectDojo API client
│   └── config/               # Configuration management
├── examples/
│   ├── go-library/           # In-process usage examples
│   └── subprocess/           # Stdio transport examples
├── Makefile                  # Build automation
└── README.md                 # This file
```

## 🔧 Configuration

### Environment Variables (for standalone binary)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DEFECTDOJO_URL` | DefectDojo instance URL | `http://localhost:8080` | Yes |
| `DEFECTDOJO_API_KEY` | API authentication key | - | Yes |
| `DEFECTDOJO_API_VERSION` | API version to use | `v2` | No |

### Programmatic Configuration (for library usage)

```go
config := &mcpserver.Config{
    DefectDojo: mcpserver.DefectDojoConfig{
        BaseURL:        "https://defectdojo.example.com",
        APIKey:         "your-secret-api-key",
        APIVersion:     "v2",
        RequestTimeout: 30 * time.Second,
    },
    Server: mcpserver.ServerConfig{
        Name:         "custom-server-name",
        Version:      "1.0.0",
        Instructions: "Custom instructions for AI agents",
    },
    Logging: mcpserver.LoggingConfig{
        Level:  "info",
        Format: "json",
    },
}
```

### Using Make commands

```bash
# Build the binary
make build

# Run as subprocess (default)
make run

# Run tests
make test

# Clean build artifacts
make clean
```

## Examples

The `examples/` directory contains complete examples showing different usage patterns:

### Go Library Example
```bash
# Shows in-process and subprocess usage
cd examples/go-library && go run main.go
```

### Subprocess Example  
```bash
# Shows stdio and sidecar server patterns
cd examples/subprocess && go run main.go
```

Both examples demonstrate:
- Different transport mechanisms (stdio, in-process)
- Tool calling patterns
- Configuration management
- Error handling

## Configuration

The server uses programmatic configuration by default, making it easy to integrate into external applications without environment files.

### Programmatic Configuration (Recommended)

When using as a Go library:

```go
import "github.com/yourusername/mcp-defect-dojo/pkg/mcpserver"

cfg := mcpserver.DefaultConfig()

// DefectDojo settings
cfg.DefectDojo.BaseURL = "https://your-defectdojo.com"
cfg.DefectDojo.APIKey = "your-api-key"
cfg.DefectDojo.RequestTimeout = 30 * time.Second

// Server settings
cfg.Server.Name = "my-defectdojo-server"
cfg.Server.Version = "v1.0.0"

// Create server
server := mcpserver.NewServer(cfg)
```

### Environment Variables

The standalone server also supports environment variables:

```bash
export DEFECTDOJO_URL="https://defectdojo.company.com"
export DEFECTDOJO_API_KEY="your-api-token"
export LOG_LEVEL="debug"
./bin/mcp-server
```

### DefectDojo Configuration
- `DEFECTDOJO_URL` - DefectDojo base URL (default: http://localhost:8080)
- `DEFECTDOJO_API_KEY` - API token for authentication
- `DEFECTDOJO_TIMEOUT` - Request timeout (default: 30s)

### Server Configuration
- `MCP_SERVER_NAME` - Server name (default: mcp-defect-dojo-server)
- `MCP_SERVER_VERSION` - Server version (default: v0.1.0)

### Logging Configuration
- `LOG_LEVEL` - Log level: debug, info, warn, error (default: info)

## MCP Tools

The server provides the following MCP tools for AI agents:

- **get_defectdojo_findings** - Retrieve vulnerability findings with filtering options
- **get_finding_detail** - Get detailed information about a specific finding  
- **mark_finding_false_positive** - Mark a finding as false positive with justification
- **defectdojo_health_check** - Verify DefectDojo connectivity and status

### Tool Examples

**Get findings:**
```json
{
  "name": "get_defectdojo_findings",
  "arguments": {
    "limit": 10,
    "offset": 0,
    "active_only": true,
    "severity": "High",
    "verified": true
  }
}
```

**Get specific finding:**
```json
{
  "name": "get_finding_detail",
  "arguments": {
    "finding_id": 123
  }
}
```

**Health check:**
```json
{
  "name": "defectdojo_health_check",
  "arguments": {}
}
```

**Mark finding as false positive:**
```json
{
  "name": "mark_finding_false_positive",
  "arguments": {
    "finding_id": 123,
    "justification": "This is a test environment vulnerability that does not affect production",
    "notes": "Verified during security review - test data only"
  }
}
```

## Development

### Project Structure

```
.
├── cmd/mcp-server/           # Standalone server binary
├── pkg/mcpserver/           # Public API for go get usage
├── internal/
│   ├── config/              # Configuration management
│   ├── defectdojo/          # DefectDojo client implementation
│   └── server/              # MCP server implementation
├── examples/                # Usage examples
├── Makefile                 # Build automation
└── go.mod                   # Go dependencies
```

### Build Commands
```bash
make help            # Show available commands
make build           # Build binary
make test            # Run tests
make clean           # Clean build artifacts
make examples        # Run examples
```

### Adding New Tools

1. Add client methods to `internal/defectdojo/client.go`
2. Register new tools in `internal/server/mcp.go`
3. Add types to `pkg/types/defectdojo.go` if needed

Example:
```go
func (s *MCPServer) registerNewTool() {
    tool := &mcp.Tool{
        Name:        "new_tool",
        Description: "Description of new tool",
        InputSchema: nil,
    }

    handler := func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
        // Tool implementation
        return &mcp.CallToolResultFor[any]{
            Content: []mcp.Content{&mcp.TextContent{Text: "Tool response"}},
        }, nil
    }

    mcp.AddTool(s.mcpServer, tool, handler)
}
```

## 🔄 Releases & Versioning

This project uses **semantic versioning** and automated releases via GitHub Actions.

### 📈 Automatic Releases

Use [conventional commits](https://www.conventionalcommits.org/) for automatic version bumping:

```bash
# Patch release (1.0.0 → 1.0.1)
git commit -m "fix: resolve authentication timeout issue"

# Minor release (1.0.0 → 1.1.0)
git commit -m "feat: add new finding filtering capabilities"

# Major release (1.0.0 → 2.0.0)  
git commit -m "feat: redesign API interface

BREAKING CHANGE: API endpoint structure has changed"
```

### 📦 Getting Releases

**Go Module:**
```bash
go get github.com/brduru/mcp-defect-dojo@latest
go get github.com/brduru/mcp-defect-dojo@v1.2.3  # specific version
```

**Pre-built Binaries:**
- [Latest Release](https://github.com/brduru/mcp-defect-dojo/releases/latest)
- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)  
- Windows (amd64)

**Build Information:**
```bash
./mcp-server --version  # Shows version, commit, build date
make version           # Shows build variables
```

For detailed release process, see [RELEASE.md](RELEASE.md).

## License

MIT License - see LICENSE file for details.
