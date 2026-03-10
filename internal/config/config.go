package config

import (
	"encoding/json"
	"os"
	"fmt"
)

// Config holds the application configuration.
type Config struct {
	Proxy           string `json:"proxy"`
	OutputFile      string `json:"output_file"`
	DefaultPassword string `json:"default_password"`
	DefaultDomain   string `json:"default_domain"`
}

const (
	DefaultProxy          = ""
	DefaultOutputFile     = "results.txt"
	DefaultConfigFilename = "config.json"
	DefaultPassword       = "" // Min 12 characters
	DefaultDomainValue    = ""
)

// DefaultConfigPath returns the default path to the config file.
func DefaultConfigPath() string {
	return DefaultConfigFilename
}

// Load reads the config from a JSON file and applies environment variable overrides.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Proxy:           DefaultProxy,
		OutputFile:      DefaultOutputFile,
		DefaultPassword: DefaultPassword,
		DefaultDomain:   DefaultDomainValue,
	}

	// Try to read the file
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	// Validate password length
	if cfg.DefaultPassword != "" && len(cfg.DefaultPassword) < 12 {
		return nil, fmt.Errorf("default_password must be at least 12 characters (got %d)", len(cfg.DefaultPassword))
	}

	// Environment variable overrides
	if proxy := os.Getenv("PROXY"); proxy != "" {
		cfg.Proxy = proxy
	}

	return cfg, nil
}
