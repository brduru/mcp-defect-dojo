package defectdojo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

// Client interface for DefectDojo API operations
type Client interface {
	GetFindings(ctx context.Context, filter types.FindingsFilter) (*types.FindingsResponse, error)
	GetFindingDetail(ctx context.Context, findingID int) (*types.Finding, error)
	MarkFalsePositive(ctx context.Context, findingID int, request types.FalsePositiveRequest) (*types.FalsePositiveResponse, error)
	HealthCheck(ctx context.Context) (bool, string)
}

// HTTPClient implements the Client interface using HTTP requests
type HTTPClient struct {
	config     *config.DefectDojoConfig
	httpClient *http.Client
}

// NewHTTPClient creates a new DefectDojo HTTP client
func NewHTTPClient(cfg *config.DefectDojoConfig) *HTTPClient {
	return &HTTPClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

// GetFindings retrieves findings from DefectDojo API with filtering
func (c *HTTPClient) GetFindings(ctx context.Context, filter types.FindingsFilter) (*types.FindingsResponse, error) {
	apiURL := fmt.Sprintf("%s%s/findings/", c.config.BaseURL, c.config.GetAPIBasePath())

	// Build query parameters
	params := url.Values{}
	params.Add("limit", strconv.Itoa(filter.Limit))
	params.Add("offset", strconv.Itoa(filter.Offset))

	if filter.ActiveOnly {
		params.Add("active", "true")
	}
	if filter.Severity != "" {
		params.Add("severity", filter.Severity)
	}
	if filter.Verified != nil {
		params.Add("verified", strconv.FormatBool(*filter.Verified))
	}
	if filter.Test != nil {
		params.Add("test", strconv.Itoa(*filter.Test))
	}

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var findings types.FindingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&findings); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &findings, nil
}

// GetFindingDetail retrieves a specific finding by ID
func (c *HTTPClient) GetFindingDetail(ctx context.Context, findingID int) (*types.Finding, error) {
	apiURL := fmt.Sprintf("%s%s/findings/%d/", c.config.BaseURL, c.config.GetAPIBasePath(), findingID)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var finding types.Finding
	if err := json.NewDecoder(resp.Body).Decode(&finding); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &finding, nil
}

// MarkFalsePositive marks a finding as false positive with justification
func (c *HTTPClient) MarkFalsePositive(ctx context.Context, findingID int, request types.FalsePositiveRequest) (*types.FalsePositiveResponse, error) {
	apiURL := fmt.Sprintf("%s%s/findings/%d/", c.config.BaseURL, c.config.GetAPIBasePath(), findingID)

	// Prepare the request payload
	payload := map[string]interface{}{
		"false_p":       true,
		"justification": request.Justification,
	}

	// Add notes if provided
	if request.Notes != "" {
		payload["notes"] = request.Notes
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var finding types.Finding
	if err := json.NewDecoder(resp.Body).Decode(&finding); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &types.FalsePositiveResponse{
		ID:      finding.ID,
		FalseP:  finding.FalseP,
		Message: "Finding successfully marked as false positive",
	}, nil
}

// HealthCheck verifies DefectDojo connectivity
func (c *HTTPClient) HealthCheck(ctx context.Context) (bool, string) {
	apiURL := fmt.Sprintf("%s%s/", c.config.BaseURL, c.config.GetAPIBasePath())

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return false, fmt.Sprintf("Failed to create request: %v", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Connection failed to %s: %v", c.config.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, fmt.Sprintf("Successfully connected to DefectDojo at %s\nAPI Version: %s\nStatus Code: %d",
			c.config.BaseURL, c.config.APIVersion, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Sprintf("DefectDojo responded with status %d: %s", resp.StatusCode, string(body))
}

// setHeaders sets common headers for API requests
func (c *HTTPClient) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Token "+c.config.APIKey)
	}
}
