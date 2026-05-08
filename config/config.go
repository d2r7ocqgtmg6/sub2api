// Package config provides configuration management for sub2api.
// It handles loading and validation of application settings from
// environment variables and configuration files.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Host string
	Port int
	BaseURL string

	// Subscription settings
	SubURL string
	SubToken string
	CacheSeconds int

	// Output format settings
	DefaultFormat string
	AllowedFormats []string

	// Security settings
	APIKey string
	EnableCORS bool
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Host:           "0.0.0.0",
		Port:           8080,
		BaseURL:        "",
		SubURL:         "",
		SubToken:       "",
		CacheSeconds:   300,
		DefaultFormat:  "clash",
		AllowedFormats: []string{"clash", "v2ray", "surge", "quantumult"},
		APIKey:         "",
		EnableCORS:     true,
	}
}

// LoadFromEnv populates the Config fields from environment variables,
// falling back to existing (default) values when variables are not set.
func (c *Config) LoadFromEnv() error {
	if v := os.Getenv("HOST"); v != "" {
		c.Host = v
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid PORT value %q: %w", v, err)
		}
		if p < 1 || p > 65535 {
			return fmt.Errorf("PORT %d is out of valid range (1-65535)", p)
		}
		c.Port = p
	}

	if v := os.Getenv("BASE_URL"); v != "" {
		c.BaseURL = strings.TrimRight(v, "/")
	}

	if v := os.Getenv("SUB_URL"); v != "" {
		c.SubURL = v
	}

	if v := os.Getenv("SUB_TOKEN"); v != "" {
		c.SubToken = v
	}

	if v := os.Getenv("CACHE_SECONDS"); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid CACHE_SECONDS value %q: %w", v, err)
		}
		if s < 0 {
			return fmt.Errorf("CACHE_SECONDS must be non-negative, got %d", s)
		}
		c.CacheSeconds = s
	}

	if v := os.Getenv("DEFAULT_FORMAT"); v != "" {
		c.DefaultFormat = strings.ToLower(v)
	}

	if v := os.Getenv("ALLOWED_FORMATS"); v != "" {
		formats := strings.Split(v, ",")
		for i, f := range formats {
			formats[i] = strings.TrimSpace(strings.ToLower(f))
		}
		c.AllowedFormats = formats
	}

	if v := os.Getenv("API_KEY"); v != "" {
		c.APIKey = v
	}

	if v := os.Getenv("ENABLE_CORS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("invalid ENABLE_CORS value %q: %w", v, err)
		}
		c.EnableCORS = b
	}

	return nil
}

// Validate checks that required configuration values are present and valid.
func (c *Config) Validate() error {
	if c.SubURL == "" {
		return fmt.Errorf("SUB_URL is required but not set")
	}

	valid := false
	for _, f := range c.AllowedFormats {
		if f == c.DefaultFormat {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("DEFAULT_FORMAT %q is not in ALLOWED_FORMATS %v", c.DefaultFormat, c.AllowedFormats)
	}

	return nil
}

// Address returns the host:port string for the HTTP server to listen on.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
