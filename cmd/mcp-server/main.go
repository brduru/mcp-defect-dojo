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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/internal/defectdojo"
	"github.com/brduru/mcp-defect-dojo/pkg/types"
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

	// Initialize DefectDojo HTTP client with configuration
	ddClient := defectdojo.NewHTTPClient(&cfg.DefectDojo)

	// Create MCP server instance with tool capabilities enabled
	s := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(true),
	)

	// Register all DefectDojo MCP tools
	addDefectDojoTools(s, ddClient)

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
	if err := server.ServeStdio(s); err != nil {
		log.Printf("❌ MCP server error: %v", err)
		os.Exit(1)
	}

	log.Printf("✅ MCP server shutdown complete")
}

func addDefectDojoTools(s *server.MCPServer, ddClient defectdojo.Client) {
	// Health check tool
	healthTool := mcp.NewTool("defectdojo_health_check",
		mcp.WithDescription("Check if DefectDojo instance is accessible and responsive"),
	)
	s.AddTool(healthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		isHealthy, message := ddClient.HealthCheck(ctx)
		status := "❌ UNHEALTHY"
		if isHealthy {
			status = "✅ HEALTHY"
		}
		return mcp.NewToolResultText(fmt.Sprintf("DefectDojo Health Check: %s\n\n%s", status, message)), nil
	})

	// Get findings tool
	findingsTool := mcp.NewTool("get_defectdojo_findings",
		mcp.WithDescription("Retrieve vulnerability findings from DefectDojo instance with optional filtering"),
		mcp.WithNumber("limit", mcp.Description("Number of findings to retrieve (default: 10)")),
		mcp.WithNumber("offset", mcp.Description("Offset for pagination (default: 0)")),
		mcp.WithBoolean("active_only", mcp.Description("Filter only active findings (default: true)")),
		mcp.WithString("severity", mcp.Description("Filter by severity (Critical, High, Medium, Low, Info)")),
		mcp.WithNumber("test", mcp.Description("Filter by test ID")),
	)
	s.AddTool(findingsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse parameters
		filter := types.FindingsFilter{
			Limit:      request.GetInt("limit", 10),
			Offset:     request.GetInt("offset", 0),
			ActiveOnly: request.GetBool("active_only", true),
			Severity:   request.GetString("severity", ""),
		}

		if test := request.GetInt("test", 0); test != 0 {
			filter.Test = &test
		}

		// Call DefectDojo API
		response, err := ddClient.GetFindings(ctx, filter)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error retrieving findings: %v", err)), nil
		}

		// Format response
		result := fmt.Sprintf("Found %d findings (showing %d):\n\n", response.Count, len(response.Results))
		for i, finding := range response.Results {
			result += fmt.Sprintf("%d. [%s] %s (ID: %d)\n", i+1, finding.Severity, finding.Title, finding.ID)
			result += fmt.Sprintf("   Active: %t, Verified: %t, False Positive: %t\n", finding.Active, finding.Verified, finding.FalseP)
			if finding.Description != "" {
				result += fmt.Sprintf("   Description: %s\n", finding.Description)
			}
			result += "\n"
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get finding detail tool
	detailTool := mcp.NewTool("get_finding_detail",
		mcp.WithDescription("Get detailed information about a specific finding by ID"),
		mcp.WithNumber("finding_id", mcp.Required(), mcp.Description("The ID of the finding to retrieve")),
	)
	s.AddTool(detailTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		findingID, err := request.RequireInt("finding_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid finding_id: %v", err)), nil
		}

		finding, err := ddClient.GetFindingDetail(ctx, findingID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error retrieving finding %d: %v", findingID, err)), nil
		}

		result := fmt.Sprintf("Finding Details (ID: %d):\n\n", finding.ID)
		result += fmt.Sprintf("Title: %s\n", finding.Title)
		result += fmt.Sprintf("Severity: %s\n", finding.Severity)
		result += fmt.Sprintf("Active: %t\n", finding.Active)
		result += fmt.Sprintf("Verified: %t\n", finding.Verified)
		result += fmt.Sprintf("False Positive: %t\n", finding.FalseP)
		result += fmt.Sprintf("Test ID: %d\n", finding.Test)
		if finding.Created != "" {
			result += fmt.Sprintf("Created: %s\n", finding.Created)
		}
		if finding.Modified != "" {
			result += fmt.Sprintf("Modified: %s\n", finding.Modified)
		}
		if finding.Description != "" {
			result += fmt.Sprintf("\nDescription:\n%s\n", finding.Description)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Mark false positive tool
	falsePositiveTool := mcp.NewTool("mark_finding_false_positive",
		mcp.WithDescription("Mark a finding as false positive with justification and optional notes/comments"),
		mcp.WithNumber("finding_id", mcp.Required(), mcp.Description("The ID of the finding to mark as false positive")),
		mcp.WithString("justification", mcp.Required(), mcp.Description("Justification for marking as false positive")),
		mcp.WithString("notes", mcp.Description("Optional additional notes or comments")),
	)
	s.AddTool(falsePositiveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		findingID, err := request.RequireInt("finding_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid finding_id: %v", err)), nil
		}

		justification, err := request.RequireString("justification")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid justification: %v", err)), nil
		}

		notes := request.GetString("notes", "")

		fpRequest := types.FalsePositiveRequest{
			IsFalsePositive: true,
			Justification:   justification,
			Notes:           notes,
		}

		response, err := ddClient.MarkFalsePositive(ctx, findingID, fpRequest)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marking finding %d as false positive: %v", findingID, err)), nil
		}

		result := fmt.Sprintf("Successfully marked finding %d as false positive:\n\n", response.ID)
		result += fmt.Sprintf("False Positive: %t\n", response.FalseP)
		result += fmt.Sprintf("Justification: %s\n", response.Justification)
		if response.Notes != "" {
			result += fmt.Sprintf("Notes: %s\n", response.Notes)
		}
		if response.Message != "" {
			result += fmt.Sprintf("Message: %s\n", response.Message)
		}

		return mcp.NewToolResultText(result), nil
	})
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
