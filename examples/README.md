# MCP DefectDojo Examples

This directory contains a complete example demonstrating the MCP DefectDojo library usage.

## 📁 Files

- **[`example.go`](example.go)** - Complete example showing both subprocess (sidecar) and embedded (in-process) usage patterns

## 🚀 Running the Example

```bash
# First, build the MCP server binary
make build

# Then run the example from the examples directory
cd examples/
go run example.go

# Or run from project root with explicit path
go run examples/example.go
```

**Note**: The sidecar example requires the MCP server binary to be built first (`make build`).

## ⚙️ Configuration

The example requires DefectDojo connection details:

```bash
# Set environment variables (optional for testing)
export DEFECTDOJO_URL="https://your-defectdojo.com"
export DEFECTDOJO_API_KEY="your-api-key"
```

**Note**: The example works even without a real DefectDojo instance - it will show expected connection errors but demonstrates that the API is working correctly.

## 📚 What the Example Demonstrates

1. **Subprocess Usage (Sidecar)**:
   - How to use the `mcp-server` binary as a separate process
   - Communication via stdio transport
   - Configuration through environment variables

2. **Embedded Usage (In-Process)**:
   - How to integrate directly into your Go code
   - In-process client creation
   - Programmatic configuration

3. **Features Tested**:
   - DefectDojo health check
   - Finding retrieval with filters
   - Error handling

## 🔧 Expected Output

```
🚀 MCP DefectDojo Complete Example
=================================

📡 Example 1: Sidecar Usage (Subprocess)
  🔄 Starting MCP server as subprocess...
  ✅ Connected to sidecar MCP server
  ✅ Health check result: [results...]

🔧 Example 2: Embedded Server Usage (In-Process)  
  🔄 Creating embedded MCP server...
  ✅ Created embedded MCP server
  ✅ Created in-process client
  ✅ Health check result: [results...]
```
