package types

import "testing"

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
