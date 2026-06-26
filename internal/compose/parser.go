package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// candidateFiles is the ordered probe list for auto-detection.
var candidateFiles = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

// Locate returns the absolute path to the compose file.
// If override is non-empty it is validated; otherwise the CWD is probed.
func Locate(override string) (string, error) {
	if override != "" {
		if _, err := os.Stat(override); err != nil {
			return "", fmt.Errorf("compose file not found: %s", override)
		}
		return filepath.Abs(override)
	}
	for _, name := range candidateFiles {
		if _, err := os.Stat(name); err == nil {
			return filepath.Abs(name)
		}
	}
	return "", fmt.Errorf("no docker-compose file found in the current directory")
}

// Parse reads and unmarshals the compose file at path into a Config.
func Parse(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	if cfg.Services == nil {
		return nil, fmt.Errorf("compose file has no services")
	}
	return &cfg, nil
}
