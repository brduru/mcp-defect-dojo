package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DefectDojo.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected default BaseURL to be http://localhost:8080, got %s", cfg.DefectDojo.BaseURL)
	}

	if cfg.DefectDojo.APIVersion != "v2" {
		t.Errorf("Expected default API version to be v2, got %s", cfg.DefectDojo.APIVersion)
	}

	if cfg.Server.Name != "mcp-defect-dojo-server" {
		t.Errorf("Expected default server name to be mcp-defect-dojo-server, got %s", cfg.Server.Name)
	}

	// Test that timeout is reasonable
	if cfg.DefectDojo.RequestTimeout <= 0 {
		t.Error("RequestTimeout should be positive")
	}
	if cfg.DefectDojo.RequestTimeout > 5*time.Minute {
		t.Error("RequestTimeout should be reasonable (< 5 minutes)")
	}
}

func TestGetAPIBasePath(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		expected   string
	}{
		{"v2 API version", "v2", "/api/v2"},
		{"v1 API version", "v1", "/api/v1"},
		{"v3 API version", "v3", "/api/v3"},
		{"empty API version", "", "/api/v2"}, // Should default to v2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &DefectDojoConfig{APIVersion: tt.apiVersion}
			result := cfg.GetAPIBasePath()
			if result != tt.expected {
				t.Errorf("GetAPIBasePath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsDebugMode(t *testing.T) {
	tests := []struct {
		level    string
		expected bool
	}{
		{"debug", true},
		{"info", false},
		{"warn", false},
		{"error", false},
		{"DEBUG", false}, // case sensitive
		{"", false},      // empty should not be debug
	}

	for _, test := range tests {
		t.Run("level_"+test.level, func(t *testing.T) {
			cfg := &LoggingConfig{Level: test.level}
			if cfg.IsDebugMode() != test.expected {
				t.Errorf("IsDebugMode() with level %s = %v, expected %v", test.level, cfg.IsDebugMode(), test.expected)
			}
		})
	}
}

// TestLoadWithEnvironment tests the configuration loading with environment variables
func TestLoadWithEnvironment(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"DEFECTDOJO_URL", "DEFECTDOJO_API_KEY", "LOG_LEVEL"}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
		os.Unsetenv(env)
	}

	// Restore environment after test
	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	t.Run("default configuration", func(t *testing.T) {
		cfg := Load()
		if cfg.DefectDojo.BaseURL != "http://localhost:8080" {
			t.Errorf("Expected default BaseURL 'http://localhost:8080', got %q", cfg.DefectDojo.BaseURL)
		}
	})

	t.Run("environment variable overrides", func(t *testing.T) {
		os.Setenv("DEFECTDOJO_URL", "https://custom.defectdojo.com")
		os.Setenv("DEFECTDOJO_API_KEY", "custom-api-key")
		os.Setenv("LOG_LEVEL", "debug")

		cfg := Load()

		if cfg.DefectDojo.BaseURL != "https://custom.defectdojo.com" {
			t.Errorf("Expected BaseURL 'https://custom.defectdojo.com', got %q", cfg.DefectDojo.BaseURL)
		}
		if cfg.DefectDojo.APIKey != "custom-api-key" {
			t.Errorf("Expected APIKey 'custom-api-key', got %q", cfg.DefectDojo.APIKey)
		}
		if cfg.Logging.Level != "debug" {
			t.Errorf("Expected log level 'debug', got %q", cfg.Logging.Level)
		}
	})
}

// BenchmarkConfigLoad benchmarks the configuration loading
func BenchmarkConfigLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Load()
	}
}

// BenchmarkConfigDefault benchmarks the default configuration creation
func BenchmarkConfigDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultConfig()
	}
}
