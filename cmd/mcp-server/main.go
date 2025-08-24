// Package main provides the standalone DefectDojo MCP server binary.
//
// This server communicates via stdio (standard input/output) for subprocess usage,
// making it compatible with MCP clients that spawn server processes.
//
// Configuration is done via environment variables for DefectDojo connection:
//   - DEFECTDOJO_URL: DefectDojo instance URL (default: http://localhost:8080)
//   - DEFECTDOJO_API_KEY: DefectDojo API token for authentication
//   - DEFECTDOJO_API_VERSION: API version to use (default: v2)
//   - LOG_LEVEL: Logging level - debug, info, warn, error (default: info)
//
// Server identity (name, version, instructions) is fixed and cannot be overridden.
//
// Example usage:
//
//	export DEFECTDOJO_URL="https://defectdojo.company.com"
//	export DEFECTDOJO_API_KEY="your-api-token"
//	./mcp-defect-dojo-server
//
// The server provides MCP tools for DefectDojo integration:
//   - defectdojo_health_check: Test DefectDojo connectivity
//   - get_defectdojo_findings: Query vulnerability findings
//   - get_finding_detail: Get detailed finding information
//   - mark_finding_false_positive: Mark findings as false positives
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/pkg/mcpserver"
)

// Version information - set at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Parse command line flags
	var showVersion = flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("mcp-defect-dojo %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		os.Exit(0)
	}

	// Setup logging to stderr since MCP protocol uses stdout for communication
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration from YAML file with environment variable overrides
	cfg := config.Load()

	// Convert to mcpserver.Config format
	mcpConfig := &mcpserver.Config{
		DefectDojo: mcpserver.DefectDojoConfig{
			BaseURL:        cfg.DefectDojo.BaseURL,
			APIKey:         cfg.DefectDojo.APIKey,
			APIVersion:     cfg.DefectDojo.APIVersion,
			RequestTimeout: cfg.DefectDojo.RequestTimeout,
		},
		Server: mcpserver.ServerConfig{
			Name:         cfg.Server.Name,
			Version:      cfg.Server.Version,
			Instructions: cfg.Server.Instructions,
		},
		Logging: mcpserver.LoggingConfig{
			Level:  cfg.Logging.Level,
			Format: cfg.Logging.Format,
		},
	}

	// Create MCP server instance
	server := mcpserver.NewServer(mcpConfig)

	// Log startup information to stderr (stdout is reserved for MCP protocol)
	log.Printf("🚀 Starting %s %s", cfg.Server.Name, cfg.Server.Version)
	log.Printf("🔗 DefectDojo URL: %s", cfg.DefectDojo.BaseURL)
	if cfg.DefectDojo.APIKey != "" {
		log.Printf("🔑 Using API key authentication")
	} else {
		log.Printf("⚠️  No API key configured - using anonymous access")
	}
	log.Printf("📡 MCP server ready for stdio communication")

	// Start the stdio server
	if err := server.Run(context.Background()); err != nil {
		log.Printf("❌ MCP server error: %v", err)
		os.Exit(1)
	}

	log.Printf("✅ MCP server shutdown complete")
}

// getEnvWithDefault retrieves an environment variable value or returns a default value.
// This utility function helps with configuration management by providing fallback values
// for optional environment variables.
//
// Parameters:
//   - key: The environment variable name to look up
//   - defaultValue: The value to return if the environment variable is not set or empty
//
// Returns:
//   - string: The environment variable value or the default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
