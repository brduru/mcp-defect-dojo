# MCP DefectDojo Examples

Examples demonstrating how to use the MCP DefectDojo server in different programming languages.

## Examples

### Go Examples
- **`subprocess/`** - Go client using MCP SDK to communicate with server subprocess
- **`go-library/`** - Go example showing direct library usage

### Other Languages
- **`python-client/`** - Python client using subprocess communication
- **`java-client/`** - Java client using subprocess communication

## Quick Start

```bash
# Build the server
make build

# Test all examples
make examples
```

## Usage Patterns

### 1. Go Library (Direct Import)
```go
import mcpdefectdojo "github.com/brduru/mcp-defect-dojo"

cfg := mcpdefectdojo.DefaultConfig()
cfg.DefectDojo.BaseURL = "https://your-defectdojo.com"
cfg.DefectDojo.APIKey = "your-api-key"

server := mcpdefectdojo.NewServer(cfg)
server.Run(ctx) // Runs as MCP subprocess
```

### 2. Subprocess (Any Language)
1. Start binary: `./bin/mcp-server`
2. Communicate via JSON-RPC over stdin/stdout
3. Use MCP protocol messages

## Available Tools

All examples demonstrate these MCP tools:
- `defectdojo_health_check` - Check DefectDojo connectivity
- `get_defectdojo_findings` - Retrieve vulnerability findings  
- `get_finding_detail` - Get specific finding details
- `mark_finding_false_positive` - Mark findings as false positive
