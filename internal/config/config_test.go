package config

import (
	"testing"
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
}

func TestGetAPIBasePath(t *testing.T) {
	cfg := &DefectDojoConfig{APIVersion: "v2"}
	expected := "/api/v2"

	if cfg.GetAPIBasePath() != expected {
		t.Errorf("Expected API base path to be %s, got %s", expected, cfg.GetAPIBasePath())
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
	}

	for _, test := range tests {
		cfg := &LoggingConfig{Level: test.level}
		if cfg.IsDebugMode() != test.expected {
			t.Errorf("IsDebugMode() with level %s = %v, expected %v", test.level, cfg.IsDebugMode(), test.expected)
		}
	}
}
