package defectdojo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brduru/mcp-defect-dojo/internal/config"
	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

func TestNewHTTPClient(t *testing.T) {
	cfg := &config.DefectDojoConfig{
		BaseURL:        "https://test.defectdojo.com",
		APIKey:         "test-api-key",
		APIVersion:     "v2",
		RequestTimeout: 30 * time.Second,
	}

	client := NewHTTPClient(cfg)

	if client == nil {
		t.Fatal("NewHTTPClient returned nil")
	}

	// Test that the client can be used (basic functionality test)
	// Since fields are private, we test via interface methods
	ctx := context.Background()
	_, _ = client.HealthCheck(ctx) // Should not panic
}

func TestHTTPClient_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedHealth bool
		expectedMsg    string
		expectError    bool
	}{
		{
			name: "healthy server",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if !strings.Contains(r.URL.Path, "/api/v2/") {
					t.Errorf("Expected API v2 path, got %s", r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "ok",
				})
			},
			expectedHealth: true,
			expectedMsg:    "",
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedHealth: false,
			expectedMsg:    "",
		},
		{
			name: "unauthorized",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			expectedHealth: false,
			expectedMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			cfg := &config.DefectDojoConfig{
				BaseURL:        server.URL,
				APIKey:         "test-key",
				APIVersion:     "v2",
				RequestTimeout: 5 * time.Second,
			}

			client := NewHTTPClient(cfg)
			healthy, msg := client.HealthCheck(context.Background())

			if healthy != tt.expectedHealth {
				t.Errorf("Expected health %v, got %v", tt.expectedHealth, healthy)
			}

			if tt.expectedMsg != "" && !strings.Contains(msg, tt.expectedMsg) {
				t.Errorf("Expected message to contain %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}

func TestHTTPClient_GetFindings(t *testing.T) {
	tests := []struct {
		name           string
		filter         types.FindingsFilter
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedCount  int
		expectError    bool
	}{
		{
			name: "successful findings request",
			filter: types.FindingsFilter{
				Limit:      10,
				ActiveOnly: true,
				Severity:   "High",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				// Check query parameters
				query := r.URL.Query()
				if query.Get("limit") != "10" {
					t.Errorf("Expected limit=10, got %s", query.Get("limit"))
				}
				if query.Get("active") != "true" {
					t.Errorf("Expected active=true, got %s", query.Get("active"))
				}
				if query.Get("severity") != "High" {
					t.Errorf("Expected severity=High, got %s", query.Get("severity"))
				}

				response := types.FindingsResponse{
					Count: 2,
					Results: []types.Finding{
						{
							ID:       1,
							Title:    "Test Finding 1",
							Severity: "High",
							Active:   true,
						},
						{
							ID:       2,
							Title:    "Test Finding 2",
							Severity: "High",
							Active:   true,
						},
					},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:   "server error",
			filter: types.FindingsFilter{Limit: 5},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectError: true,
		},
		{
			name:   "invalid JSON response",
			filter: types.FindingsFilter{Limit: 5},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("invalid json"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			cfg := &config.DefectDojoConfig{
				BaseURL:        server.URL,
				APIKey:         "test-key",
				APIVersion:     "v2",
				RequestTimeout: 5 * time.Second,
			}

			client := NewHTTPClient(cfg)
			response, err := client.GetFindings(context.Background(), tt.filter)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Fatal("Response is nil")
			}

			if len(response.Results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(response.Results))
			}
		})
	}
}

func TestHTTPClient_GetFindingDetail(t *testing.T) {
	tests := []struct {
		name           string
		findingID      int
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedTitle  string
		expectError    bool
	}{
		{
			name:      "successful finding detail request",
			findingID: 123,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				if !strings.Contains(r.URL.Path, "/123/") {
					t.Errorf("Expected finding ID 123 in path, got %s", r.URL.Path)
				}

				finding := types.Finding{
					ID:          123,
					Title:       "Detailed Test Finding",
					Severity:    "Critical",
					Active:      true,
					Verified:    true,
					Description: "This is a detailed test finding",
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(finding)
			},
			expectedTitle: "Detailed Test Finding",
			expectError:   false,
		},
		{
			name:      "finding not found",
			findingID: 999,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			cfg := &config.DefectDojoConfig{
				BaseURL:        server.URL,
				APIKey:         "test-key",
				APIVersion:     "v2",
				RequestTimeout: 5 * time.Second,
			}

			client := NewHTTPClient(cfg)
			finding, err := client.GetFindingDetail(context.Background(), tt.findingID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if finding == nil {
				t.Fatal("Finding is nil")
			}

			if finding.Title != tt.expectedTitle {
				t.Errorf("Expected title %q, got %q", tt.expectedTitle, finding.Title)
			}
		})
	}
}

func TestHTTPClient_MarkFalsePositive(t *testing.T) {
	tests := []struct {
		name           string
		findingID      int
		request        types.FalsePositiveRequest
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectError    bool
	}{
		{
			name:      "successful false positive marking",
			findingID: 456,
			request: types.FalsePositiveRequest{
				IsFalsePositive: true,
				Justification:   "This is a test environment",
				Notes:           "Confirmed with security team",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Expected PATCH request, got %s", r.Method)
				}

				// Verify request body
				var reqBody types.FalsePositiveRequest
				json.NewDecoder(r.Body).Decode(&reqBody)

				if !reqBody.IsFalsePositive {
					t.Error("Expected IsFalsePositive to be true")
				}

				response := types.FalsePositiveResponse{
					ID:            456,
					FalseP:        true,
					Justification: reqBody.Justification,
					Notes:         reqBody.Notes,
					Message:       "Successfully marked as false positive",
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectError: false,
		},
		{
			name:      "server error",
			findingID: 456,
			request: types.FalsePositiveRequest{
				IsFalsePositive: true,
				Justification:   "Test",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			cfg := &config.DefectDojoConfig{
				BaseURL:        server.URL,
				APIKey:         "test-key",
				APIVersion:     "v2",
				RequestTimeout: 5 * time.Second,
			}

			client := NewHTTPClient(cfg)
			response, err := client.MarkFalsePositive(context.Background(), tt.findingID, tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Fatal("Response is nil")
			}

			if response.ID != tt.findingID {
				t.Errorf("Expected ID %d, got %d", tt.findingID, response.ID)
			}
		})
	}
}

func TestHTTPClient_ContextCancellation(t *testing.T) {
	// Test that context cancellation is properly handled
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.DefectDojoConfig{
		BaseURL:        server.URL,
		APIKey:         "test-key",
		APIVersion:     "v2",
		RequestTimeout: 5 * time.Second,
	}

	client := NewHTTPClient(cfg)

	// Create a context that will be cancelled quickly
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.GetFindings(ctx, types.FindingsFilter{Limit: 10})
	if err == nil {
		t.Error("Expected context cancellation error but got none")
	}
}

func TestHTTPClient_AuthenticationHeaders(t *testing.T) {
	expectedAPIKey := "test-api-key-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Token " + expectedAPIKey

		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuth, authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types.FindingsResponse{Count: 0, Results: []types.Finding{}})
	}))
	defer server.Close()

	cfg := &config.DefectDojoConfig{
		BaseURL:        server.URL,
		APIKey:         expectedAPIKey,
		APIVersion:     "v2",
		RequestTimeout: 5 * time.Second,
	}

	client := NewHTTPClient(cfg)
	_, err := client.GetFindings(context.Background(), types.FindingsFilter{Limit: 1})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
