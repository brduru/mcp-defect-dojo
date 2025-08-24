package mcpserver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brduru/mcp-defect-dojo/pkg/types"
)

// MockDefectDojoClient implements the defectdojo.Client interface for testing
type MockDefectDojoClient struct {
	HealthCheckFunc       func(ctx context.Context) (bool, string)
	GetFindingsFunc       func(ctx context.Context, filter types.FindingsFilter) (*types.FindingsResponse, error)
	GetFindingDetailFunc  func(ctx context.Context, findingID int) (*types.Finding, error)
	MarkFalsePositiveFunc func(ctx context.Context, findingID int, request types.FalsePositiveRequest) (*types.FalsePositiveResponse, error)
}

func (m *MockDefectDojoClient) HealthCheck(ctx context.Context) (bool, string) {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return true, "Mock DefectDojo is healthy"
}

func (m *MockDefectDojoClient) GetFindings(ctx context.Context, filter types.FindingsFilter) (*types.FindingsResponse, error) {
	if m.GetFindingsFunc != nil {
		return m.GetFindingsFunc(ctx, filter)
	}
	return &types.FindingsResponse{
		Count: 2,
		Results: []types.Finding{
			{
				ID:          1,
				Title:       "Test Finding 1",
				Severity:    "High",
				Active:      true,
				Verified:    false,
				Description: "Test description 1",
				Test:        100,
			},
			{
				ID:          2,
				Title:       "Test Finding 2",
				Severity:    "Medium",
				Active:      true,
				Verified:    true,
				Description: "Test description 2",
				Test:        101,
			},
		},
	}, nil
}

func (m *MockDefectDojoClient) GetFindingDetail(ctx context.Context, findingID int) (*types.Finding, error) {
	if m.GetFindingDetailFunc != nil {
		return m.GetFindingDetailFunc(ctx, findingID)
	}
	if findingID == 999 {
		return nil, fmt.Errorf("finding not found: %d", findingID)
	}
	return &types.Finding{
		ID:          findingID,
		Title:       fmt.Sprintf("Test Finding %d", findingID),
		Severity:    "High",
		Active:      true,
		Verified:    false,
		Description: fmt.Sprintf("Detailed description for finding %d", findingID),
		Test:        100,
		Created:     "2023-01-01T00:00:00Z",
		Modified:    "2023-01-02T00:00:00Z",
	}, nil
}

func (m *MockDefectDojoClient) MarkFalsePositive(ctx context.Context, findingID int, request types.FalsePositiveRequest) (*types.FalsePositiveResponse, error) {
	if m.MarkFalsePositiveFunc != nil {
		return m.MarkFalsePositiveFunc(ctx, findingID, request)
	}
	if findingID == 999 {
		return nil, fmt.Errorf("finding not found: %d", findingID)
	}
	return &types.FalsePositiveResponse{
		ID:            findingID,
		FalseP:        true,
		Justification: request.Justification,
		Notes:         request.Notes,
		Message:       "Successfully marked as false positive",
	}, nil
}

// Test configuration creation and validation
func TestNewServer(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   bool // true if server should be created successfully
	}{
		{
			name: "valid configuration",
			config: &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL:        "https://defectdojo.example.com",
					APIKey:         "test-api-key",
					APIVersion:     "v2",
					RequestTimeout: 30 * time.Second,
				},
				Server: ServerConfig{
					Name:         "test-server",
					Version:      "1.0.0",
					Instructions: "Test instructions",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "text",
				},
			},
			want: true,
		},
		{
			name: "minimal configuration",
			config: &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL: "http://localhost:8080",
					APIKey:  "test-key",
				},
				Server: ServerConfig{
					Name:    "minimal-server",
					Version: "1.0.0",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.config)
			if (server != nil) != tt.want {
				t.Errorf("NewServer() = %v, want %v", server != nil, tt.want)
			}
			if server != nil {
				if server.mcpServer == nil {
					t.Error("NewServer() created server with nil mcpServer")
				}
				if server.ddClient == nil {
					t.Error("NewServer() created server with nil ddClient")
				}
			}
		})
	}
}

// Test server lifecycle methods
func TestServerLifecycle(t *testing.T) {
	config := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL:    "http://localhost:8080",
			APIKey:     "test-key",
			APIVersion: "v2",
		},
		Server: ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	server := NewServer(config)
	if server == nil {
		t.Fatal("Failed to create server")
	}

	// Test GetMCPServer
	mcpServer := server.GetMCPServer()
	if mcpServer == nil {
		t.Error("GetMCPServer() returned nil")
	}

	// Test that we can get the same instance
	mcpServer2 := server.GetMCPServer()
	if mcpServer != mcpServer2 {
		t.Error("GetMCPServer() returned different instances")
	}
}

// Test NewServerWithAPIKey
func TestNewServerWithAPIKey(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		want   bool
	}{
		{
			name:   "valid API key",
			apiKey: "test-api-key-123",
			want:   true,
		},
		{
			name:   "empty API key",
			apiKey: "",
			want:   true, // Should still create server, just without auth
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServerWithAPIKey(tt.apiKey)
			if (err == nil && server != nil) != tt.want {
				t.Errorf("NewServerWithAPIKey() error = %v, server = %v, want %v", err, server != nil, tt.want)
			}
		})
	}
}

// Test NewServerWithSettings
func TestNewServerWithSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings DefectDojoSettings
		want     bool
	}{
		{
			name: "full settings",
			settings: DefectDojoSettings{
				BaseURL:    "https://defectdojo.company.com",
				APIKey:     "test-api-key",
				APIVersion: "v2",
			},
			want: true,
		},
		{
			name: "minimal settings",
			settings: DefectDojoSettings{
				BaseURL: "http://localhost:8080",
				APIKey:  "test-key",
			},
			want: true,
		},
		{
			name:     "empty settings",
			settings: DefectDojoSettings{},
			want:     true, // Should use defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServerWithSettings(tt.settings)
			if (err == nil && server != nil) != tt.want {
				t.Errorf("NewServerWithSettings() error = %v, server = %v, want %v", err, server != nil, tt.want)
			}
		})
	}
}

// Test configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name            string
		config          *Config
		expectNilServer bool
	}{
		{
			name:            "nil config uses defaults",
			config:          nil,
			expectNilServer: false, // Should create server with defaults
		},
		{
			name: "valid config",
			config: &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL: "http://localhost:8080",
					APIKey:  "test-key",
				},
				Server: ServerConfig{
					Name:    "test-server",
					Version: "1.0.0",
				},
			},
			expectNilServer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.config)
			isNil := server == nil

			if isNil != tt.expectNilServer {
				t.Errorf("NewServer() returned nil = %v, expected nil = %v", isNil, tt.expectNilServer)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewServer(b *testing.B) {
	config := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL:    "http://localhost:8080",
			APIKey:     "test-key",
			APIVersion: "v2",
		},
		Server: ServerConfig{
			Name:    "benchmark-server",
			Version: "1.0.0",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := NewServer(config)
		_ = server.GetMCPServer()
	}
}

func BenchmarkNewServerWithAPIKey(b *testing.B) {
	apiKey := "benchmark-api-key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server, err := NewServerWithAPIKey(apiKey)
		if err != nil {
			b.Fatalf("NewServerWithAPIKey() error = %v", err)
		}
		_ = server.GetMCPServer()
	}
}

// Test server creation with various timeout configurations
func TestServerTimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{
			name:    "default timeout",
			timeout: 0,                // Should use default
			want:    30 * time.Second, // Assuming this is the default
		},
		{
			name:    "custom timeout",
			timeout: 60 * time.Second,
			want:    60 * time.Second,
		},
		{
			name:    "very short timeout",
			timeout: 1 * time.Second,
			want:    1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL:        "http://localhost:8080",
					APIKey:         "test-key",
					RequestTimeout: tt.timeout,
				},
				Server: ServerConfig{
					Name:    "timeout-test-server",
					Version: "1.0.0",
				},
			}

			server := NewServer(config)
			if server == nil {
				t.Fatal("Failed to create server")
			}

			// Test that server was created successfully
			mcpServer := server.GetMCPServer()
			if mcpServer == nil {
				t.Error("GetMCPServer() returned nil")
			}
		})
	}
}

// Test that server properly handles different API versions
func TestServerAPIVersions(t *testing.T) {
	versions := []string{"v1", "v2", "v3", ""}

	for _, version := range versions {
		t.Run(fmt.Sprintf("version_%s", version), func(t *testing.T) {
			config := &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL:    "http://localhost:8080",
					APIKey:     "test-key",
					APIVersion: version,
				},
				Server: ServerConfig{
					Name:    "version-test-server",
					Version: "1.0.0",
				},
			}

			server := NewServer(config)
			if server == nil {
				t.Fatal("Failed to create server with API version:", version)
			}

			mcpServer := server.GetMCPServer()
			if mcpServer == nil {
				t.Error("GetMCPServer() returned nil for API version:", version)
			}
		})
	}
}

func TestDefectDojoSettingsValidation(t *testing.T) {
	tests := []struct {
		name        string
		settings    DefectDojoSettings
		expectedURL string
	}{
		{
			name: "HTTPS URL",
			settings: DefectDojoSettings{
				BaseURL: "https://secure.defectdojo.com",
				APIKey:  "secure-key",
			},
			expectedURL: "https://secure.defectdojo.com",
		},
		{
			name: "HTTP URL for dev",
			settings: DefectDojoSettings{
				BaseURL: "http://localhost:8080",
				APIKey:  "dev-key",
			},
			expectedURL: "http://localhost:8080",
		},
		{
			name: "URL with path",
			settings: DefectDojoSettings{
				BaseURL: "https://company.com/defectdojo",
				APIKey:  "company-key",
			},
			expectedURL: "https://company.com/defectdojo",
		},
		{
			name: "minimal settings",
			settings: DefectDojoSettings{
				BaseURL: "https://minimal.com",
				APIKey:  "key",
			},
			expectedURL: "https://minimal.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServerWithSettings(tt.settings)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if server == nil {
				t.Error("Expected server but got nil")
			}
		})
	}
}

func TestConfigurationChain(t *testing.T) {
	// Test the configuration chain from settings to final config
	settings := DefectDojoSettings{
		BaseURL:    "https://test.example.com",
		APIKey:     "test-config-key",
		APIVersion: "v2",
	}

	server, err := NewServerWithSettings(settings)
	if err != nil {
		t.Fatalf("Failed to create server with settings: %v", err)
	}

	if server == nil {
		t.Fatal("Server is nil")
	}

	mcpServer := server.GetMCPServer()
	if mcpServer == nil {
		t.Error("MCP server is nil")
	}
}

func TestServerCreationMethods(t *testing.T) {
	t.Run("NewServerWithAPIKey creates valid server", func(t *testing.T) {
		server, err := NewServerWithAPIKey("test-api-key-123")
		if err != nil {
			t.Errorf("Failed to create server: %v", err)
		}
		if server == nil {
			t.Error("Server is nil")
		}
		if server.GetMCPServer() == nil {
			t.Error("MCP server is nil")
		}
	})

	t.Run("NewServerWithSettings creates valid server", func(t *testing.T) {
		settings := DefectDojoSettings{
			BaseURL:    "https://test.defectdojo.com",
			APIKey:     "settings-test-key",
			APIVersion: "v2",
		}

		server, err := NewServerWithSettings(settings)
		if err != nil {
			t.Errorf("Failed to create server: %v", err)
		}
		if server == nil {
			t.Error("Server is nil")
		}
		if server.GetMCPServer() == nil {
			t.Error("MCP server is nil")
		}
	})
}

func TestTimeoutConfigurations(t *testing.T) {
	timeouts := []time.Duration{
		10 * time.Second,
		30 * time.Second,
		60 * time.Second,
		2 * time.Minute,
	}

	for _, timeout := range timeouts {
		t.Run(fmt.Sprintf("timeout_%s", timeout), func(t *testing.T) {
			config := &Config{
				DefectDojo: DefectDojoConfig{
					BaseURL:        "https://timeout.test.com",
					APIKey:         "timeout-test-key",
					RequestTimeout: timeout,
				},
			}

			server := NewServer(config)
			if server == nil {
				t.Errorf("Failed to create server with timeout %s", timeout)
			}
		})
	}
}

func TestEmptyAndNilConfigurations(t *testing.T) {
	t.Run("empty DefectDojoSettings", func(t *testing.T) {
		settings := DefectDojoSettings{}
		server, err := NewServerWithSettings(settings)
		if err != nil {
			t.Errorf("Unexpected error with empty settings: %v", err)
		}
		if server == nil {
			t.Error("Server should not be nil with empty settings")
		}
	})

	t.Run("empty API key", func(t *testing.T) {
		server, err := NewServerWithAPIKey("")
		if err != nil {
			t.Errorf("Unexpected error with empty API key: %v", err)
		}
		if server == nil {
			t.Error("Server should not be nil with empty API key")
		}
	})
}

func TestMCPToolsIntegration(t *testing.T) {
	cfg := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL: "https://test.defectdojo.com",
			APIKey:  "test-key",
		},
		Server: ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("Failed to create server")
	}

	mcpServer := server.GetMCPServer()
	if mcpServer == nil {
		t.Fatal("MCP server is nil")
	}

	// Test that tools are registered (indirect test via server creation)
	t.Run("server_creation_with_tools", func(t *testing.T) {
		// The fact that NewServer completes without error indicates
		// that addDefectDojoTools ran successfully
		if server == nil {
			t.Error("Expected server to be created with tools")
		}
	})
}

func TestServerRunMethodExists(t *testing.T) {
	// Test that the Run method exists and can be called
	// We can't easily test the full stdio functionality in unit tests
	// but we can test that the method exists and doesn't panic immediately
	cfg := &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL: "https://test.com",
			APIKey:  "test-key",
		},
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("Failed to create server")
	}

	// Test that Run method exists and can be called
	// Note: In a real scenario, this would start stdio communication
	// For testing, we just verify the method exists and is callable
	t.Run("run_method_exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Run method caused panic: %v", r)
			}
		}()

		// We can't actually test stdio communication in unit tests
		// but we can test that the method can be called
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		// This will likely timeout or fail, but should not panic
		_ = server.Run(ctx)
	})
}

func TestConfigurationEdgeCases(t *testing.T) {
	t.Run("nil_config_uses_defaults", func(t *testing.T) {
		// Test that NewServer handles nil config by using defaults
		server := NewServer(nil)
		if server == nil {
			t.Error("Expected server to be created with default config, but got nil")
		}
	})

	t.Run("config_with_empty_values", func(t *testing.T) {
		cfg := &Config{
			DefectDojo: DefectDojoConfig{
				BaseURL: "",
				APIKey:  "",
			},
			Server: ServerConfig{
				Name:    "",
				Version: "",
			},
		}

		server := NewServer(cfg)
		if server == nil {
			t.Error("Expected server to be created even with empty config values")
		}
	})
}

func TestAPIVersionHandling(t *testing.T) {
	versions := []string{"v1", "v2", "v3", ""}

	for _, version := range versions {
		t.Run("api_version_"+version, func(t *testing.T) {
			settings := DefectDojoSettings{
				BaseURL:    "https://test.example.com",
				APIKey:     "test-key",
				APIVersion: version,
			}

			server, err := NewServerWithSettings(settings)
			if err != nil {
				t.Errorf("Failed to create server with API version %s: %v", version, err)
			}
			if server == nil {
				t.Errorf("Server is nil for API version %s", version)
			}
		})
	}
}
