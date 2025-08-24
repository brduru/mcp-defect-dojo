package types

import (
	"encoding/json"
	"testing"
)

// TestFindingsFilter tests the FindingsFilter structure and its methods
func TestFindingsFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter FindingsFilter
	}{
		{
			name: "default filter",
			filter: FindingsFilter{
				Limit:      10,
				Offset:     0,
				ActiveOnly: true,
			},
		},
		{
			name: "filter with all fields",
			filter: FindingsFilter{
				Limit:      20,
				Offset:     10,
				ActiveOnly: false,
				Severity:   "High",
				Verified:   boolPtr(true),
				Test:       intPtr(123),
			},
		},
		{
			name: "filter with pointers nil",
			filter: FindingsFilter{
				Limit:      5,
				Offset:     0,
				ActiveOnly: true,
				Severity:   "Medium",
				Verified:   nil,
				Test:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the filter can be created and fields are accessible
			if tt.filter.Limit <= 0 {
				t.Error("Limit should be positive")
			}
			if tt.filter.Offset < 0 {
				t.Error("Offset should be non-negative")
			}
		})
	}
}

// TestFinding tests the Finding structure
func TestFinding(t *testing.T) {
	finding := Finding{
		ID:          1,
		Title:       "Test SQL Injection",
		Severity:    "High",
		Active:      true,
		Verified:    false,
		FalseP:      false,
		Description: "A SQL injection vulnerability was found",
		Test:        100,
		Created:     "2023-01-01T00:00:00Z",
		Modified:    "2023-01-02T00:00:00Z",
	}

	// Test JSON marshaling and unmarshaling
	data, err := json.Marshal(finding)
	if err != nil {
		t.Fatalf("Failed to marshal finding: %v", err)
	}

	var unmarshaled Finding
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal finding: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaled.ID != finding.ID {
		t.Errorf("ID mismatch: got %d, want %d", unmarshaled.ID, finding.ID)
	}
	if unmarshaled.Title != finding.Title {
		t.Errorf("Title mismatch: got %q, want %q", unmarshaled.Title, finding.Title)
	}
	if unmarshaled.Severity != finding.Severity {
		t.Errorf("Severity mismatch: got %q, want %q", unmarshaled.Severity, finding.Severity)
	}
	if unmarshaled.Active != finding.Active {
		t.Errorf("Active mismatch: got %t, want %t", unmarshaled.Active, finding.Active)
	}
	if unmarshaled.Verified != finding.Verified {
		t.Errorf("Verified mismatch: got %t, want %t", unmarshaled.Verified, finding.Verified)
	}
	if unmarshaled.FalseP != finding.FalseP {
		t.Errorf("FalseP mismatch: got %t, want %t", unmarshaled.FalseP, finding.FalseP)
	}
}

// TestFindingsResponse tests the FindingsResponse structure
func TestFindingsResponse(t *testing.T) {
	response := FindingsResponse{
		Count: 2,
		Next:  stringPtr("https://api.example.com/findings/?offset=20"),
		Previous: nil,
		Results: []Finding{
			{
				ID:       1,
				Title:    "Finding 1",
				Severity: "High",
				Active:   true,
			},
			{
				ID:       2,
				Title:    "Finding 2",
				Severity: "Medium",
				Active:   false,
			},
		},
	}

	// Test that count matches results length
	if response.Count != len(response.Results) {
		t.Errorf("Count mismatch: got %d, want %d", response.Count, len(response.Results))
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled FindingsResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Count != response.Count {
		t.Errorf("Count mismatch after unmarshal: got %d, want %d", unmarshaled.Count, response.Count)
	}
}

// TestFalsePositiveRequest tests the FalsePositiveRequest structure
func TestFalsePositiveRequest(t *testing.T) {
	tests := []struct {
		name    string
		request FalsePositiveRequest
		valid   bool
	}{
		{
			name: "valid request with justification",
			request: FalsePositiveRequest{
				IsFalsePositive: true,
				Justification:   "This is a test environment",
				Notes:           "Additional notes here",
			},
			valid: true,
		},
		{
			name: "valid request without notes",
			request: FalsePositiveRequest{
				IsFalsePositive: true,
				Justification:   "False positive due to configuration",
				Notes:           "",
			},
			valid: true,
		},
		{
			name: "request with false positive false",
			request: FalsePositiveRequest{
				IsFalsePositive: false,
				Justification:   "Not marking as false positive",
				Notes:           "",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var unmarshaled FalsePositiveRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify fields are preserved
			if unmarshaled.IsFalsePositive != tt.request.IsFalsePositive {
				t.Errorf("IsFalsePositive mismatch: got %t, want %t", 
					unmarshaled.IsFalsePositive, tt.request.IsFalsePositive)
			}
			if unmarshaled.Justification != tt.request.Justification {
				t.Errorf("Justification mismatch: got %q, want %q", 
					unmarshaled.Justification, tt.request.Justification)
			}
			if unmarshaled.Notes != tt.request.Notes {
				t.Errorf("Notes mismatch: got %q, want %q", 
					unmarshaled.Notes, tt.request.Notes)
			}
		})
	}
}

// TestFalsePositiveResponse tests the FalsePositiveResponse structure
func TestFalsePositiveResponse(t *testing.T) {
	response := FalsePositiveResponse{
		ID:            123,
		FalseP:        true,
		Justification: "Test justification",
		Notes:         "Test notes",
		Message:       "Successfully marked as false positive",
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled FalsePositiveResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify all fields
	if unmarshaled.ID != response.ID {
		t.Errorf("ID mismatch: got %d, want %d", unmarshaled.ID, response.ID)
	}
	if unmarshaled.FalseP != response.FalseP {
		t.Errorf("FalseP mismatch: got %t, want %t", unmarshaled.FalseP, response.FalseP)
	}
	if unmarshaled.Justification != response.Justification {
		t.Errorf("Justification mismatch: got %q, want %q", unmarshaled.Justification, response.Justification)
	}
}

// Test edge cases for JSON handling
func TestJSONEdgeCases(t *testing.T) {
	t.Run("empty finding", func(t *testing.T) {
		finding := Finding{}
		data, err := json.Marshal(finding)
		if err != nil {
			t.Fatalf("Failed to marshal empty finding: %v", err)
		}

		var unmarshaled Finding
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal empty finding: %v", err)
		}
	})

	t.Run("finding with special characters", func(t *testing.T) {
		finding := Finding{
			Title:       "Test with special chars: <script>alert('xss')</script>",
			Description: "Description with unicode: 测试 🔒 ñoño",
			Severity:    "High",
		}

		data, err := json.Marshal(finding)
		if err != nil {
			t.Fatalf("Failed to marshal finding with special chars: %v", err)
		}

		var unmarshaled Finding
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal finding with special chars: %v", err)
		}

		if unmarshaled.Title != finding.Title {
			t.Errorf("Title with special chars not preserved: got %q, want %q", 
				unmarshaled.Title, finding.Title)
		}
		if unmarshaled.Description != finding.Description {
			t.Errorf("Description with special chars not preserved: got %q, want %q", 
				unmarshaled.Description, finding.Description)
		}
	})
}

// Benchmark tests
func BenchmarkFindingMarshal(b *testing.B) {
	finding := Finding{
		ID:          1,
		Title:       "Benchmark Test Finding",
		Severity:    "High",
		Active:      true,
		Verified:    false,
		Description: "This is a benchmark test finding with a longer description to test performance",
		Test:        100,
		Created:     "2023-01-01T00:00:00Z",
		Modified:    "2023-01-02T00:00:00Z",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(finding)
		if err != nil {
			b.Fatalf("Marshal error: %v", err)
		}
	}
}

func BenchmarkFindingUnmarshal(b *testing.B) {
	finding := Finding{
		ID:          1,
		Title:       "Benchmark Test Finding",
		Severity:    "High",
		Active:      true,
		Verified:    false,
		Description: "This is a benchmark test finding with a longer description to test performance",
		Test:        100,
		Created:     "2023-01-01T00:00:00Z",
		Modified:    "2023-01-02T00:00:00Z",
	}

	data, err := json.Marshal(finding)
	if err != nil {
		b.Fatalf("Setup error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var unmarshaled Finding
		err := json.Unmarshal(data, &unmarshaled)
		if err != nil {
			b.Fatalf("Unmarshal error: %v", err)
		}
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

// TestIsValidSeverity tests the severity validation function
func TestIsValidSeverity(t *testing.T) {
	tests := []struct {
		severity string
		expected bool
	}{
		{"Critical", true},
		{"High", true},
		{"Medium", true},
		{"Low", true},
		{"Info", true},
		{"Unknown", false},
		{"", false},
		{"CRITICAL", false}, // case sensitive
	}

	for _, test := range tests {
		result := IsValidSeverity(test.severity)
		if result != test.expected {
			t.Errorf("IsValidSeverity(%q) = %v, expected %v", test.severity, result, test.expected)
		}
	}
}

// TestValidSeverities tests the function that returns all valid severities
func TestValidSeverities(t *testing.T) {
	severities := ValidSeverities()
	expected := []string{"Info", "Low", "Medium", "High", "Critical"}

	if len(severities) != len(expected) {
		t.Fatalf("Expected %d severities, got %d", len(expected), len(severities))
	}

	for i, severity := range expected {
		if severities[i] != severity {
			t.Errorf("Expected severity[%d] = %q, got %q", i, severity, severities[i])
		}
	}
}
