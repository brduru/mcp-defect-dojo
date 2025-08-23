package main

import (
	"context"
	"log"
	"time"

	"github.com/brduru/mcp-defect-dojo/pkg/mcpserver"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	log.Printf("🚀 MCP DefectDojo Go Library Example")
	log.Printf("====================================")

	// Example 1: Using as a sidecar (subprocess)
	log.Printf("\n📡 Example 1: Sidecar Usage (Subprocess)")
	if err := runSidecarExample(); err != nil {
		log.Printf("❌ Sidecar example failed: %v", err)
	}

	// Example 2: Embedding the server (in-process)
	log.Printf("\n🔧 Example 2: Embedded Server Usage (In-Process)")
	if err := runEmbeddedExample(); err != nil {
		log.Printf("❌ Embedded example failed: %v", err)
	}
}

func runSidecarExample() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("  🔄 Starting MCP server as subprocess...")

	// Create MCP client for stdio transport (sidecar)
	mcpClient, err := client.NewStdioMCPClient("../../bin/mcp-server", []string{
		"DEFECTDOJO_URL=http://localhost:8080",
		"DEFECTDOJO_API_KEY=your-api-key-here", // Replace with your actual API key
	})
	if err != nil {
		return err
	}
	defer mcpClient.Close()

	log.Printf("  ✅ Connected to sidecar MCP server")

	// Initialize the client
	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "mcp-defectdojo-sidecar-example",
				Version: "1.0.0",
			},
		},
	}

	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Printf("  ⚠️ Client initialization error: %v", err)
		return err
	}

	// Test health check tool
	result, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "defectdojo_health_check",
			Arguments: map[string]any{},
		},
	})
	if err != nil {
		log.Printf("  ❌ Health check failed: %v", err)
		return err
	}

	log.Printf("  ✅ Health check result: %v", result.Content)

	// Test get findings with filter
	result, err = mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_defectdojo_findings",
			Arguments: map[string]any{
				"limit":       3,
				"active_only": true,
			},
		},
	})
	if err != nil {
		log.Printf("  ❌ Get findings failed: %v", err)
		return err
	}

	log.Printf("  ✅ Found findings: %v", result.Content)
	return nil
}

func runEmbeddedExample() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Printf("  🔄 Creating embedded MCP server...")

	// Option 1: Simple usage with just API key (uses localhost:8080)
	// server, err := mcpserver.NewServerWithAPIKey("your-api-key-here")

	// Option 2: Full control with custom DefectDojo settings
	server, err := mcpserver.NewServerWithSettings(mcpserver.DefectDojoSettings{
		BaseURL:    "http://localhost:8080", // Your DefectDojo URL
		APIKey:     "your-api-key-here",     // Your API key
		APIVersion: "v2",                    // API version (optional, defaults to v2)
	})
	if err != nil {
		return err
	}

	log.Printf("  ✅ Created embedded MCP server")

	// Create in-process client
	mcpClient, err := client.NewInProcessClient(server.GetMCPServer())
	if err != nil {
		return err
	}

	log.Printf("  ✅ Created in-process client")

	// Initialize the client
	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "mcp-defectdojo-example",
				Version: "1.0.0",
			},
		},
	}

	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Printf("  ⚠️ Client initialization error: %v", err)
		return err
	}

	// Test health check tool
	result, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "defectdojo_health_check",
			Arguments: map[string]any{},
		},
	})
	if err != nil {
		log.Printf("  ❌ Health check failed: %v", err)
		return err
	}

	log.Printf("  ✅ Health check result: %v", result.Content)

	// Test get findings with filter
	result, err = mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_defectdojo_findings",
			Arguments: map[string]any{
				"limit":       3,
				"active_only": true,
				"severity":    "High",
			},
		},
	})
	if err != nil {
		log.Printf("  ❌ Get findings failed: %v", err)
		return err
	}

	log.Printf("  ✅ Found findings: %v", result.Content)
	return nil
}
