package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VirusTotalAPIKey string   `yaml:"virustotal_api_key"`
	AlienVaultAPIKey string   `yaml:"alienvault_api_key"`
	// Additional configuration options:
	Concurrency int      `yaml:"concurrency"`       // Number of concurrent requests.
	UserAgents  []string `yaml:"user_agents"`       // Custom list of user-agent strings.
	RateLimit   int      `yaml:"rate_limit"`        // Optional rate limit (requests per minute).
	// You can add more fields as needed.
}

// LoadConfig reads a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = "config.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate checks that required configuration fields are set.
func Validate(cfg *Config) error {
	if cfg.VirusTotalAPIKey == "" {
		return errors.New("virus total api key is missing")
	}
	if cfg.AlienVaultAPIKey == "" {
		return errors.New("alien vault api key is missing")
	}
	if cfg.Concurrency <= 0 {
		// Set a default value if not provided.
		cfg.Concurrency = 5
	}
	if len(cfg.UserAgents) == 0 {
		// Set default user agents.
		cfg.UserAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)",
		}
	}
	return nil
}
