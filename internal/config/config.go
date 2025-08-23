package config

import (
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	DefectDojo DefectDojoConfig
	Server     ServerConfig
	Logging    LoggingConfig
}

// DefectDojoConfig contains DefectDojo API configuration
type DefectDojoConfig struct {
	BaseURL        string
	APIKey         string
	APIVersion     string
	RequestTimeout time.Duration
}

// ServerConfig contains MCP server configuration
type ServerConfig struct {
	Name         string
	Version      string
	Instructions string
	Host         string
	Port         int
	Transport    string // "stdio", "http"
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		DefectDojo: DefectDojoConfig{
			BaseURL:        "http://localhost:8080",
			APIKey:         "",
			APIVersion:     "v2",
			RequestTimeout: 30 * time.Second,
		},
		Server: ServerConfig{
			Name:         "mcp-defect-dojo-server",
			Version:      "0.2.1",
			Instructions: "MCP server for DefectDojo integration. Provides tools to query vulnerability findings and manage security data.",
			Host:         "localhost",
			Port:         8000,
			Transport:    "stdio", // Default to stdio for subprocess usage
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// GetTimeout returns request timeout as string for http client
func (c *DefectDojoConfig) GetTimeout() string {
	return c.RequestTimeout.String()
}

// GetAPIBasePath returns the full API base path
func (c *DefectDojoConfig) GetAPIBasePath() string {
	return "/api/" + c.APIVersion
}

// IsDebugMode checks if debug logging is enabled
func (c *LoggingConfig) IsDebugMode() bool {
	return c.Level == "debug"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// Load loads configuration with defaults and environment variable overrides
// DefectDojo settings can be overridden, but server identity remains fixed
func Load() *Config {
	// Start with default configuration (fixed server identity)
	config := DefaultConfig()

	// Override ONLY DefectDojo settings with environment variables
	if val := os.Getenv("DEFECTDOJO_URL"); val != "" {
		config.DefectDojo.BaseURL = val
	}
	if val := os.Getenv("DEFECTDOJO_API_KEY"); val != "" {
		config.DefectDojo.APIKey = val
	}
	if val := os.Getenv("DEFECTDOJO_API_VERSION"); val != "" {
		config.DefectDojo.APIVersion = val
	}

	// Logging can be overridden for debugging
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.Logging.Level = val
	}
	if val := os.Getenv("LOG_FORMAT"); val != "" {
		config.Logging.Format = val
	}

	// Server identity (name, version, instructions) should NOT be overrideable
	// These are part of the library's identity and should remain consistent

	return config
}
