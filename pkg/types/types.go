package types

// Finding represents a DefectDojo finding/vulnerability with all core fields.
// This structure mirrors the DefectDojo API response for individual findings.
//
// Example:
//
//	finding := &Finding{
//		ID:          123,
//		Title:       "SQL Injection in Login",
//		Severity:    "Critical",
//		Description: "SQL injection vulnerability found in authentication endpoint",
//		Active:      true,
//		Verified:    true,
//		FalseP:      false,
//	}
type Finding struct {
	ID          int    `json:"id"`                   // Unique finding identifier
	Title       string `json:"title"`               // Finding title/summary
	Severity    string `json:"severity"`            // Severity level (Critical, High, Medium, Low, Info)
	Description string `json:"description"`         // Detailed finding description
	Active      bool   `json:"active"`              // Whether the finding is currently active
	Verified    bool   `json:"verified"`            // Whether the finding has been verified
	FalseP      bool   `json:"false_p"`             // Whether marked as false positive
	Test        int    `json:"test"`                // Associated test ID
	Created     string `json:"created,omitempty"`   // Creation timestamp (ISO 8601)
	Modified    string `json:"modified,omitempty"`  // Last modification timestamp (ISO 8601)
}

// FalsePositiveRequest represents a request to mark a finding as false positive.
// This structure is used when updating a finding's false positive status via the API.
//
// Example:
//
//	request := &FalsePositiveRequest{
//		IsFalsePositive: true,
//		Justification:   "This is expected behavior in test environment",
//		Notes:          "Confirmed with security team",
//	}
type FalsePositiveRequest struct {
	IsFalsePositive bool   `json:"false_p"`           // Whether to mark as false positive
	Justification   string `json:"justification,omitempty"` // Reason for marking as false positive
	Notes           string `json:"notes,omitempty"`   // Additional notes or comments
}

// FalsePositiveResponse represents the response from marking a finding as false positive.
// This structure contains the updated finding information after the false positive operation.
type FalsePositiveResponse struct {
	ID            int    `json:"id"`                   // Finding ID that was updated
	FalseP        bool   `json:"false_p"`              // Updated false positive status
	Justification string `json:"justification,omitempty"` // Applied justification
	Notes         string `json:"notes,omitempty"`      // Applied notes
	Message       string `json:"message,omitempty"`    // Optional response message from API
}

// FindingsResponse represents the paginated API response for findings list queries.
// This follows DefectDojo's standard pagination format for bulk finding retrieval.
//
// Example usage:
//
//	response := &FindingsResponse{}
//	// After API call...
//	for _, finding := range response.Results {
//		fmt.Printf("Finding %d: %s\n", finding.ID, finding.Title)
//	}
type FindingsResponse struct {
	Count    int       `json:"count"`     // Total number of findings matching the query
	Next     *string   `json:"next"`      // URL for next page of results (nil if last page)
	Previous *string   `json:"previous"`  // URL for previous page of results (nil if first page)
	Results  []Finding `json:"results"`   // Array of findings for current page
}

// FindingsFilter contains filtering and pagination options for findings queries.
// Use this structure to control which findings are returned and how they're paginated.
//
// Example:
//
//	filter := &FindingsFilter{
//		Limit:      50,              // Return up to 50 results
//		ActiveOnly: true,            // Only active findings
//		Severity:   "Critical",      // Only critical severity
//		Verified:   &[]bool{true}[0], // Only verified findings
//		Offset:     0,               // Start from beginning
//	}
type FindingsFilter struct {
	Limit      int    // Maximum number of results to return (default: 100)
	ActiveOnly bool   // Filter to only active findings
	Severity   string // Filter by severity level (Critical, High, Medium, Low, Info)
	Verified   *bool  // Filter by verification status (nil = all, true = verified only, false = unverified only)
	Test       *int   // Filter by specific test ID (nil = all tests)
	Offset     int    // Number of results to skip for pagination
}

// Severity level constants for DefectDojo findings.
// These constants represent the standard severity levels used in DefectDojo
// vulnerability management. Use these constants instead of string literals
// to avoid typos and ensure consistency.
const (
	SeverityInfo     = "Info"     // Informational findings (lowest severity)
	SeverityLow      = "Low"      // Low severity vulnerabilities  
	SeverityMedium   = "Medium"   // Medium severity vulnerabilities
	SeverityHigh     = "High"     // High severity vulnerabilities
	SeverityCritical = "Critical" // Critical severity vulnerabilities (highest severity)
)

// ValidSeverities returns a slice of all valid severity levels in DefectDojo.
// This function is useful for validation, UI dropdowns, and documentation.
//
// Returns severity levels in ascending order of criticality:
// ["Info", "Low", "Medium", "High", "Critical"]
//
// Example:
//
//	severities := ValidSeverities()
//	for _, severity := range severities {
//		fmt.Printf("Valid severity: %s\n", severity)
//	}
func ValidSeverities() []string {
	return []string{
		SeverityInfo,
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}
}

// IsValidSeverity checks if the provided severity level is valid in DefectDojo.
// This function performs case-sensitive validation against the standard severity levels.
//
// Parameters:
//   - severity: The severity string to validate
//
// Returns:
//   - true if the severity is valid (Info, Low, Medium, High, or Critical)
//   - false if the severity is invalid or empty
//
// Example:
//
//	if IsValidSeverity("Critical") {
//		fmt.Println("Valid severity level")
//	}
//	
//	if !IsValidSeverity("invalid") {
//		fmt.Println("Invalid severity level")
//	}
func IsValidSeverity(severity string) bool {
	for _, valid := range ValidSeverities() {
		if severity == valid {
			return true
		}
	}
	return false
}
