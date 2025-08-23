package types

// Finding represents a DefectDojo finding/vulnerability
type Finding struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Verified    bool   `json:"verified"`
	FalseP      bool   `json:"false_p"`
	Test        int    `json:"test"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

// FalsePositiveRequest represents a request to mark finding as false positive
type FalsePositiveRequest struct {
	IsFalsePositive bool   `json:"false_p"`
	Justification   string `json:"justification,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

// FalsePositiveResponse represents the response from marking false positive
type FalsePositiveResponse struct {
	ID            int    `json:"id"`
	FalseP        bool   `json:"false_p"`
	Justification string `json:"justification,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Message       string `json:"message,omitempty"`
}

// FindingsResponse represents the API response for findings list
type FindingsResponse struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Finding `json:"results"`
}

// FindingsFilter contains filtering options for findings queries
type FindingsFilter struct {
	Limit      int
	ActiveOnly bool
	Severity   string
	Verified   *bool
	Test       *int
	Offset     int
}

// Severity levels as constants
const (
	SeverityInfo     = "Info"
	SeverityLow      = "Low"
	SeverityMedium   = "Medium"
	SeverityHigh     = "High"
	SeverityCritical = "Critical"
)

// ValidSeverities returns a slice of valid severity levels
func ValidSeverities() []string {
	return []string{
		SeverityInfo,
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}
}

// IsValidSeverity checks if a severity level is valid
func IsValidSeverity(severity string) bool {
	for _, valid := range ValidSeverities() {
		if severity == valid {
			return true
		}
	}
	return false
}
