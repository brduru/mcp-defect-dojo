package defectdojo

import (
	"context"
	"fmt"

	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

// MockClient implements Client interface for testing
type MockClient struct {
	findings      *types.FindingsResponse
	finding       *types.Finding
	healthStatus  bool
	healthMessage string
	shouldError   bool
}

// NewMockClient creates a new mock client for testing
func NewMockClient() *MockClient {
	return &MockClient{
		findings: &types.FindingsResponse{
			Count: 2,
			Results: []types.Finding{
				{
					ID:          1,
					Title:       "SQL Injection vulnerability",
					Severity:    types.SeverityHigh,
					Description: "SQL injection found in login form",
					Active:      true,
					Verified:    true,
					Test:        100,
				},
				{
					ID:          2,
					Title:       "XSS vulnerability",
					Severity:    types.SeverityMedium,
					Description: "Cross-site scripting vulnerability",
					Active:      true,
					Verified:    false,
					Test:        101,
				},
			},
		},
		finding: &types.Finding{
			ID:          1,
			Title:       "SQL Injection vulnerability",
			Severity:    types.SeverityHigh,
			Description: "Detailed description of SQL injection vulnerability",
			Active:      true,
			Verified:    true,
			Test:        100,
		},
		healthStatus:  true,
		healthMessage: "Mock DefectDojo is healthy",
		shouldError:   false,
	}
}

// GetFindings returns mock findings
func (m *MockClient) GetFindings(ctx context.Context, filter types.FindingsFilter) (*types.FindingsResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.findings, nil
}

// GetFindingDetail returns mock finding detail
func (m *MockClient) GetFindingDetail(ctx context.Context, findingID int) (*types.Finding, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.finding, nil
}

// MarkFalsePositive mocks marking a finding as false positive
func (m *MockClient) MarkFalsePositive(ctx context.Context, findingID int, request types.FalsePositiveRequest) (*types.FalsePositiveResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return &types.FalsePositiveResponse{
		ID:            findingID,
		FalseP:        true,
		Justification: request.Justification,
		Notes:         request.Notes,
		Message:       "Mock: Finding marked as false positive",
	}, nil
}

// HealthCheck returns mock health status
func (m *MockClient) HealthCheck(ctx context.Context) (bool, string) {
	return m.healthStatus, m.healthMessage
}

// SetError configures the mock to return errors
func (m *MockClient) SetError(shouldError bool) {
	m.shouldError = shouldError
}

// SetHealthStatus configures mock health status
func (m *MockClient) SetHealthStatus(healthy bool, message string) {
	m.healthStatus = healthy
	m.healthMessage = message
}
