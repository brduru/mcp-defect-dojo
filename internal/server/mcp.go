// Package server provides the internal MCP server implementation for DefectDojo integration.
// This package contains the core server logic and tool registration functionality.
package server

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/internal/defectdojo"
	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

// MCPServer wraps the MCP server functionality for DefectDojo integration.
// It manages the configuration, DefectDojo client, and registered MCP tools.
type MCPServer struct {
	config           *config.Config    // Server configuration
	defectDojoClient defectdojo.Client // DefectDojo API client
	mcpServer        *server.MCPServer // Underlying MCP server from mcp-go
}

// NewMCPServer creates a new MCP server instance with DefectDojo integration.
// It initializes the DefectDojo client and registers all available tools.
//
// Parameters:
//   - cfg: Configuration containing DefectDojo connection settings and server metadata
//
// Returns:
//   - *MCPServer: A configured MCP server ready to handle DefectDojo operations
func NewMCPServer(cfg *config.Config) *MCPServer {
	client := defectdojo.NewHTTPClient(&cfg.DefectDojo)

	mcpServer := server.NewMCPServer("mcp-defect-dojo", "1.0.0",
		server.WithToolCapabilities(true),
	)

	s := &MCPServer{
		config:           cfg,
		defectDojoClient: client,
		mcpServer:        mcpServer,
	}

	s.registerTools()
	return s
}

// GetServer returns the underlying MCP server for in-process usage.
// This enables direct integration with MCP clients in the same process.
//
// Returns:
//   - *server.MCPServer: The mcp-go server instance for in-process communication
func (s *MCPServer) GetServer() *server.MCPServer {
	return s.mcpServer
}

// ServeStdio starts the server with stdio transport for subprocess communication.
// This method blocks until the server is terminated or an error occurs.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error that occurs during server operation
func (s *MCPServer) ServeStdio(ctx context.Context) error {
	return server.ServeStdio(s.mcpServer)
}

// registerTools registers all available DefectDojo MCP tools with the server.
// This includes health check, findings retrieval, finding details, and false positive marking.
func (s *MCPServer) registerTools() {
	s.registerGetFindingsTool()
	s.registerGetFindingDetailTool()
	s.registerMarkFalsePositiveTool()
	s.registerHealthCheckTool()
}

// registerGetFindingsTool registers the get_defectdojo_findings tool.
// This tool allows querying DefectDojo for vulnerability findings with various filters.
func (s *MCPServer) registerGetFindingsTool() {
	tool := mcp.Tool{
		Name:        "get_defectdojo_findings",
		Description: "Retrieve vulnerability findings from DefectDojo instance with optional filtering",
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("🔍 Tool call: get_defectdojo_findings with params: %+v", request.Params.Arguments)

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			args = make(map[string]any)
		}

		filter := s.parseFilterFromParams(args)
		log.Printf("📊 Parsed filter: %+v", filter)

		findings, err := s.defectDojoClient.GetFindings(ctx, filter)
		if err != nil {
			log.Printf("❌ Error retrieving findings: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Error retrieving findings: %v", err)), nil
		}

		log.Printf("✅ Retrieved %d findings successfully", len(findings.Results))
		result := s.formatFindingsResponse(findings)
		return mcp.NewToolResultText(result), nil
	}

	s.mcpServer.AddTool(tool, handler)
}

// registerGetFindingDetailTool registers the get_finding_detail tool.
// This tool retrieves detailed information about a specific finding by ID.
func (s *MCPServer) registerGetFindingDetailTool() {
	tool := mcp.Tool{
		Name:        "get_finding_detail",
		Description: "Get detailed information about a specific finding by ID",
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("🔍 Tool call: get_finding_detail with params: %+v", request.Params.Arguments)

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			args = make(map[string]any)
		}

		findingID := s.extractFindingID(args)
		if findingID == 0 {
			log.Printf("❌ Invalid or missing finding ID in params")
			return mcp.NewToolResultError("Error: finding_id parameter is required and must be a positive integer"), nil
		}

		log.Printf("📋 Getting details for finding ID: %d", findingID)
		finding, err := s.defectDojoClient.GetFindingDetail(ctx, findingID)
		if err != nil {
			log.Printf("❌ Error retrieving finding detail: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Error retrieving finding %d: %v", findingID, err)), nil
		}

		log.Printf("✅ Retrieved finding detail successfully: ID %d, Title: %s", finding.ID, finding.Title)
		result := s.formatFindingDetail(finding)
		return mcp.NewToolResultText(result), nil
	}

	s.mcpServer.AddTool(tool, handler)
}

// registerMarkFalsePositiveTool registers the mark_finding_false_positive tool.
// This tool allows marking findings as false positives with justification.
func (s *MCPServer) registerMarkFalsePositiveTool() {
	tool := mcp.Tool{
		Name:        "mark_finding_false_positive",
		Description: "Mark a finding as false positive with optional justification",
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("🔍 Tool call: mark_finding_false_positive with params: %+v", request.Params.Arguments)

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Error: Invalid arguments"), nil
		}

		findingID := s.extractFindingID(args)
		if findingID == 0 {
			return mcp.NewToolResultError("Error: finding_id parameter is required and must be a positive integer"), nil
		}

		log.Printf("🔄 Marking finding %d as false positive", findingID)
		// Create a false positive request with optional justification
		fpRequest := types.FalsePositiveRequest{
			IsFalsePositive: true,
			Justification:   "Marked as false positive via MCP tool",
		}

		// Add justification if provided
		if justification, ok := args["justification"].(string); ok && justification != "" {
			fpRequest.Justification = justification
		}

		// Add notes if provided
		if notes, ok := args["notes"].(string); ok && notes != "" {
			fpRequest.Notes = notes
		}

		_, err := s.defectDojoClient.MarkFalsePositive(ctx, findingID, fpRequest)
		if err != nil {
			log.Printf("❌ Error marking finding as false positive: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Error marking finding %d as false positive: %v", findingID, err)), nil
		}

		log.Printf("✅ Successfully marked finding %d as false positive", findingID)
		result := fmt.Sprintf("Successfully marked finding %d as false positive", findingID)
		if fpRequest.Justification != "" {
			result += fmt.Sprintf("\nJustification: %s", fpRequest.Justification)
		}
		return mcp.NewToolResultText(result), nil
	}

	s.mcpServer.AddTool(tool, handler)
}

// registerHealthCheckTool registers the defectdojo_health_check tool.
// This tool verifies DefectDojo connectivity and instance health.
func (s *MCPServer) registerHealthCheckTool() {
	tool := mcp.Tool{
		Name:        "defectdojo_health_check",
		Description: "Check DefectDojo instance health and connectivity",
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("🔍 Tool call: defectdojo_health_check")

		ok, message := s.defectDojoClient.HealthCheck(ctx)
		if !ok {
			log.Printf("❌ DefectDojo health check failed: %s", message)
			return mcp.NewToolResultError(fmt.Sprintf("DefectDojo health check failed: %s", message)), nil
		}

		log.Printf("✅ DefectDojo health check passed")
		return mcp.NewToolResultText("DefectDojo health check passed successfully"), nil
	}

	s.mcpServer.AddTool(tool, handler)
}

// parseFilterFromParams extracts and validates filter parameters from MCP tool arguments.
// It provides sensible defaults for optional parameters.
//
// Parameters:
//   - arguments: Map of parameters passed from the MCP client
//
// Returns:
//   - types.FindingsFilter: Validated filter parameters with defaults applied
func (s *MCPServer) parseFilterFromParams(arguments map[string]any) types.FindingsFilter {
	filter := types.FindingsFilter{
		Limit:      20,   // default
		ActiveOnly: true, // default
		Offset:     0,
	}

	if v, ok := arguments["limit"]; ok {
		if limitFloat, ok := v.(float64); ok {
			filter.Limit = int(limitFloat)
		}
	}

	if v, ok := arguments["offset"]; ok {
		if offsetFloat, ok := v.(float64); ok {
			filter.Offset = int(offsetFloat)
		}
	}

	if v, ok := arguments["active_only"]; ok {
		if activeBool, ok := v.(bool); ok {
			filter.ActiveOnly = activeBool
		}
	}

	if v, ok := arguments["severity"]; ok {
		if severityStr, ok := v.(string); ok && types.IsValidSeverity(severityStr) {
			filter.Severity = severityStr
		}
	}

	if v, ok := arguments["verified"]; ok {
		if verifiedBool, ok := v.(bool); ok {
			filter.Verified = &verifiedBool
		}
	}

	if v, ok := arguments["test"]; ok {
		if testFloat, ok := v.(float64); ok {
			testInt := int(testFloat)
			filter.Test = &testInt
		}
	}

	return filter
}

// extractFindingID extracts and validates the finding ID from MCP tool arguments.
// It handles both numeric and string representations of the ID.
//
// Parameters:
//   - arguments: Map of parameters passed from the MCP client
//
// Returns:
//   - int: The finding ID, or 0 if invalid or missing
func (s *MCPServer) extractFindingID(arguments map[string]any) int {
	if v, ok := arguments["finding_id"]; ok {
		switch val := v.(type) {
		case float64:
			if val > 0 {
				return int(val)
			}
		case string:
			if id, err := strconv.Atoi(val); err == nil && id > 0 {
				return id
			}
		case int:
			if val > 0 {
				return val
			}
		}
	}
	return 0
}

// formatFindingsResponse formats the findings response for display to the user.
// It creates a human-readable summary with key information for each finding.
//
// Parameters:
//   - findings: The findings response from DefectDojo API
//
// Returns:
//   - string: Formatted text suitable for display to AI agents
func (s *MCPServer) formatFindingsResponse(findings *types.FindingsResponse) string {
	if findings.Count == 0 {
		return "No findings found matching the specified criteria."
	}

	result := fmt.Sprintf("Found %d findings (showing first %d):\n\n",
		findings.Count, len(findings.Results))

	for i, finding := range findings.Results {
		if i >= 10 { // Limit display to first 10
			result += fmt.Sprintf("... and %d more findings\n", len(findings.Results)-10)
			break
		}

		status := "Active"
		if !finding.Active {
			status = "Inactive"
		}

		verified := ""
		if finding.Verified {
			verified = " (Verified)"
		}

		result += fmt.Sprintf("%d. [%s] %s%s\n", finding.ID, finding.Severity, finding.Title, verified)
		result += fmt.Sprintf("   Status: %s | Test: %d\n", status, finding.Test)
		result += fmt.Sprintf("   %s\n\n", finding.Description)
	}

	if findings.Next != nil {
		result += "Note: More results available. Use 'offset' parameter to paginate.\n"
	}

	return result
}

// formatFindingDetail formats a single finding for detailed display.
// It includes all available information about the finding in a structured format.
//
// Parameters:
//   - finding: The detailed finding information from DefectDojo API
//
// Returns:
//   - string: Formatted text with comprehensive finding details
func (s *MCPServer) formatFindingDetail(finding *types.Finding) string {
	result := fmt.Sprintf("Finding #%d Details:\n\n", finding.ID)
	result += fmt.Sprintf("Title: %s\n", finding.Title)
	result += fmt.Sprintf("Severity: %s\n", finding.Severity)
	result += fmt.Sprintf("Active: %t\n", finding.Active)
	result += fmt.Sprintf("Verified: %t\n", finding.Verified)
	result += fmt.Sprintf("Test ID: %d\n", finding.Test)

	if finding.Created != "" {
		result += fmt.Sprintf("Created: %s\n", finding.Created)
	}
	if finding.Modified != "" {
		result += fmt.Sprintf("Modified: %s\n", finding.Modified)
	}

	result += fmt.Sprintf("\nDescription:\n%s\n", finding.Description)

	return result
}
