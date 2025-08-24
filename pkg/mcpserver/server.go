// Package mcpserver provides Model Context Protocol (MCP) server integration for DefectDojo.
//
// This package enables AI agents to interact with DefectDojo vulnerability management platform
// through the Model Context Protocol. It supports multiple transport methods and provides
// comprehensive tools for vulnerability management workflows.
//
// # Overview
//
// The mcpserver package is the primary public API for integrating DefectDojo with MCP-compatible
// AI tools. It provides a clean, well-documented interface for both embedded usage and
// subprocess communication patterns.
//
// # Supported MCP Tools
//
//   - defectdojo_health_check: Verify DefectDojo API connectivity and health status
//   - get_defectdojo_findings: Retrieve and filter vulnerability findings with advanced options
//   - get_finding_detail: Get comprehensive details about specific vulnerabilities
//   - mark_finding_false_positive: Mark findings as false positives with audit trail
//
// # Transport Methods
//
// The server supports two primary communication patterns:
//
//   - In-Process: Direct function calls for embedded usage within Go applications
//   - Stdio: Subprocess communication for language-agnostic integration
//
// # Quick Start Examples
//
// ## Simple Setup with API Key
//
//	server, err := mcpserver.NewServerWithAPIKey("your-defectdojo-api-key")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create in-process client
//	client, err := client.NewInProcessClient(server.GetMCPServer())
//	if err != nil {
//		log.Fatal(err)
//	}
//
// ## Custom Configuration
//
//	server, err := mcpserver.NewServerWithSettings(mcpserver.DefectDojoSettings{
//		BaseURL:    "https://defectdojo.company.com",
//		APIKey:     "your-api-key",
//		APIVersion: "v2",
//	})
//
// ## Full Configuration Control
//
//	config := &mcpserver.Config{
//		DefectDojo: mcpserver.DefectDojoConfig{
//			BaseURL:        "https://defectdojo.company.com",
//			APIKey:         "your-api-key",
//			APIVersion:     "v2",
//			RequestTimeout: 30 * time.Second,
//		},
//		Server: mcpserver.ServerConfig{
//			Name:         "custom-defectdojo-server",
//			Version:      "1.0.0",
//			Instructions: "Custom instructions for AI agents",
//		},
//	}
//	server := mcpserver.NewServer(config)
//
// ## Subprocess Usage
//
//	// Start the server as subprocess
//	stdioClient, err := client.NewStdioMCPClient("./mcp-defect-dojo-server", []string{
//		"DEFECTDOJO_URL=https://your-defectdojo.com",
//		"DEFECTDOJO_API_KEY=your-api-key",
//	})
//
// # Error Handling
//
// All MCP tools properly propagate errors through the MCP protocol. Connection failures,
// authentication errors, and API errors are returned as proper Go errors that can be
// handled by MCP clients.
//
// # Thread Safety
//
// The server instances are thread-safe and can be used concurrently from multiple goroutines.
// The underlying HTTP client handles connection pooling and concurrent requests efficiently.
package mcpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/internal/defectdojo"
	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

// Server represents an MCP DefectDojo server instance
type Server struct {
	mcpServer *server.MCPServer
	ddClient  defectdojo.Client
}

// Config represents the server configuration for the DefectDojo MCP server.
// This structure contains all the necessary settings to connect to DefectDojo
// and configure the MCP server behavior.
type Config struct {
	DefectDojo DefectDojoConfig // DefectDojo API connection settings
	Server     ServerConfig     // MCP server metadata and behavior
	Logging    LoggingConfig    // Logging configuration
}

// DefectDojoConfig contains DefectDojo API configuration.
// These settings control how the server connects to and interacts with DefectDojo.
type DefectDojoConfig struct {
	BaseURL        string        // DefectDojo instance URL (e.g., "https://defectdojo.company.com")
	APIKey         string        // DefectDojo API token for authentication
	APIVersion     string        // DefectDojo API version to use (typically "v2")
	RequestTimeout time.Duration // HTTP request timeout for DefectDojo API calls
}

// ServerConfig contains MCP server configuration.
// These settings define the server's identity and behavior in the MCP protocol.
type ServerConfig struct {
	Name         string // Server name as reported to MCP clients
	Version      string // Server version for client compatibility
	Instructions string // Optional instructions displayed to AI agents
}

// LoggingConfig contains logging configuration.
// Controls how the server logs information for debugging and monitoring.
type LoggingConfig struct {
	Level  string // Log level: "debug", "info", "warn", "error"
	Format string // Log format: "text", "json"
}

// NewServer creates a new MCP DefectDojo server with the provided configuration.
// The server supports multiple transport methods: in-process and stdio.
//
// Parameters:
//   - cfg: Configuration containing DefectDojo API settings, server info, and logging options
//
// Returns:
//   - *Server: A configured MCP server ready to handle DefectDojo operations
//
// The server automatically registers the following MCP tools:
//   - get_defectdojo_findings: Query vulnerability findings with filters
//   - get_finding_detail: Get detailed information about a specific finding
//   - mark_finding_false_positive: Mark findings as false positives with justification
//   - defectdojo_health_check: Test DefectDojo API connectivity
func NewServer(cfg *Config) *Server {
	// Use default config if nil is provided
	if cfg == nil {
		defaultCfg := config.DefaultConfig()
		cfg = &Config{
			DefectDojo: DefectDojoConfig{
				BaseURL:        defaultCfg.DefectDojo.BaseURL,
				APIKey:         defaultCfg.DefectDojo.APIKey,
				APIVersion:     defaultCfg.DefectDojo.APIVersion,
				RequestTimeout: defaultCfg.DefectDojo.RequestTimeout,
			},
			Server: ServerConfig{
				Name:         defaultCfg.Server.Name,
				Version:      defaultCfg.Server.Version,
				Instructions: defaultCfg.Server.Instructions,
			},
			Logging: LoggingConfig{
				Level:  defaultCfg.Logging.Level,
				Format: defaultCfg.Logging.Format,
			},
		}
	}

	// Create DefectDojo client
	ddClient := defectdojo.NewHTTPClient(&config.DefectDojoConfig{
		BaseURL:        cfg.DefectDojo.BaseURL,
		APIKey:         cfg.DefectDojo.APIKey,
		APIVersion:     cfg.DefectDojo.APIVersion,
		RequestTimeout: cfg.DefectDojo.RequestTimeout,
	})

	// Create MCP server using mcp-go
	mcpServer := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(true),
	)

	// Add DefectDojo tools
	addDefectDojoTools(mcpServer, ddClient)

	return &Server{
		mcpServer: mcpServer,
		ddClient:  ddClient,
	}
}

// NewServerWithAPIKey creates a new MCP DefectDojo server using default configuration with API key override.
// This is a simple method for embedded usage where you only need to set the API key.
//
// Parameters:
//   - apiKey: DefectDojo API key to use
//
// Returns:
//   - *Server: A configured MCP server ready to handle DefectDojo operations
//   - error: Any error that occurs during configuration loading or server creation
func NewServerWithAPIKey(apiKey string) (*Server, error) {
	// Load configuration with defaults and environment variable overrides
	cfg := config.Load()

	// Override API key
	cfg.DefectDojo.APIKey = apiKey

	// Convert to mcpserver.Config format
	mcpConfig := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL:        cfg.DefectDojo.BaseURL,
			APIKey:         cfg.DefectDojo.APIKey,
			APIVersion:     cfg.DefectDojo.APIVersion,
			RequestTimeout: cfg.DefectDojo.RequestTimeout,
		},
		Server: ServerConfig{
			Name:         cfg.Server.Name,
			Version:      cfg.Server.Version,
			Instructions: cfg.Server.Instructions,
		},
		Logging: LoggingConfig{
			Level:  cfg.Logging.Level,
			Format: cfg.Logging.Format,
		},
	}

	return NewServer(mcpConfig), nil
}

// DefectDojoSettings contains DefectDojo connection settings for embedded usage
type DefectDojoSettings struct {
	BaseURL    string // DefectDojo instance URL (e.g., "https://defectdojo.company.com")
	APIKey     string // DefectDojo API key for authentication
	APIVersion string // DefectDojo API version (default: "v2")
}

// NewServerWithSettings creates a new MCP DefectDojo server with custom DefectDojo settings.
// This provides full control over DefectDojo connection for embedded usage.
//
// Parameters:
//   - settings: DefectDojo connection settings (URL, API key, version)
//
// Returns:
//   - *Server: A configured MCP server ready to handle DefectDojo operations
//   - error: Any error that occurs during server creation
func NewServerWithSettings(settings DefectDojoSettings) (*Server, error) {
	// Start with default configuration for server identity and logging
	cfg := config.DefaultConfig()

	// Override DefectDojo settings
	cfg.DefectDojo.BaseURL = settings.BaseURL
	cfg.DefectDojo.APIKey = settings.APIKey

	if settings.APIVersion != "" {
		cfg.DefectDojo.APIVersion = settings.APIVersion
	}

	// Convert to mcpserver.Config format
	mcpConfig := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL:        cfg.DefectDojo.BaseURL,
			APIKey:         cfg.DefectDojo.APIKey,
			APIVersion:     cfg.DefectDojo.APIVersion,
			RequestTimeout: cfg.DefectDojo.RequestTimeout,
		},
		Server: ServerConfig{
			Name:         cfg.Server.Name,
			Version:      cfg.Server.Version,
			Instructions: cfg.Server.Instructions,
		},
		Logging: LoggingConfig{
			Level:  cfg.Logging.Level,
			Format: cfg.Logging.Format,
		},
	}

	return NewServer(mcpConfig), nil
}

// Run starts the MCP server with stdio transport.
// This method is typically used for subprocess communication where the server
// communicates with a parent process via standard input/output.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error that occurs during server operation
//
// This is the primary method for subprocess/sidecar usage patterns.
func (s *Server) Run(ctx context.Context) error {
	return server.ServeStdio(s.mcpServer)
}

// GetMCPServer returns the underlying MCP server for in-process use.
// This enables direct integration with MCP clients in the same process,
// avoiding the overhead of network or stdio communication.
//
// Returns:
//   - *server.MCPServer: The mcp-go server instance for direct method calls
//
// Use this method when you want to embed the DefectDojo MCP server
// directly in your application for maximum performance and simplicity.
// This is useful for creating in-process clients with client.NewInProcessClient().
func (s *Server) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

// Available MCP Tools:
//
// The DefectDojo MCP server provides the following tools for AI agents:
//
// - defectdojo_health_check: Test connectivity to DefectDojo instance
//   Returns the health status and version information
//
// - get_defectdojo_findings: Query vulnerability findings with filters
//   Supports pagination, severity filtering, and active/inactive status
//
// - get_finding_detail: Get comprehensive details for a specific finding
//   Returns full vulnerability information including CVSS scores and descriptions
//
// - mark_finding_false_positive: Mark findings as false positives
//   Requires justification and supports additional notes for audit trail

// addDefectDojoTools registers all DefectDojo MCP tools with the server.
// This function sets up the tool handlers and their JSON schemas for parameter validation.
func addDefectDojoTools(s *server.MCPServer, ddClient defectdojo.Client) {
	// Health check tool
	healthTool := mcp.NewTool("defectdojo_health_check",
		mcp.WithDescription("Check if DefectDojo instance is accessible and responsive"),
	)
	s.AddTool(healthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		isHealthy, message := ddClient.HealthCheck(ctx)
		if !isHealthy {
			return nil, fmt.Errorf("DefectDojo Health Check failed: %s", message)
		}
		return mcp.NewToolResultText(fmt.Sprintf("DefectDojo Health Check: ✅ HEALTHY\n\n%s", message)), nil
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
			return nil, fmt.Errorf("error retrieving findings: %w", err)
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
			return nil, fmt.Errorf("invalid finding_id: %w", err)
		}

		finding, err := ddClient.GetFindingDetail(ctx, findingID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving finding %d: %w", findingID, err)
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
			return nil, fmt.Errorf("invalid finding_id: %w", err)
		}

		justification, err := request.RequireString("justification")
		if err != nil {
			return nil, fmt.Errorf("invalid justification: %w", err)
		}

		notes := request.GetString("notes", "")

		fpRequest := types.FalsePositiveRequest{
			IsFalsePositive: true,
			Justification:   justification,
			Notes:           notes,
		}

		response, err := ddClient.MarkFalsePositive(ctx, findingID, fpRequest)
		if err != nil {
			return nil, fmt.Errorf("error marking finding %d as false positive: %w", findingID, err)
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
