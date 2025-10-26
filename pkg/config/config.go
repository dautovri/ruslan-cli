package config

import (
"fmt"
"os"
"path/filepath"

"gopkg.in/yaml.v3"
)

// Environment represents a Vault environment configuration
type Environment struct {
	Name        string `yaml:"name"`
	ProjectID   string `yaml:"project_id"`
	Region      string `yaml:"region"`
	ClusterName string `yaml:"cluster_name"`
	Namespace   string `yaml:"namespace"`
	ServiceName string `yaml:"service_name"`
	VaultAddr   string `yaml:"vault_addr,omitempty"`
	VaultPort   string `yaml:"vault_port"`
	UseNipIO    bool   `yaml:"use_nipio"`
	Token       string `yaml:"token,omitempty"` // Added Token field
}

// Config represents the CLI configuration
type Config struct {
	CurrentEnvironment string                  `yaml:"current_environment"`
	Environments       map[string]*Environment `yaml:"environments"`
	TokenFile          string                  `yaml:"token_file"`
	CacheDir           string                  `yaml:"cache_dir"`
	OutputFormat       string                  `yaml:"output_format"`
	AutoRefresh        bool                    `yaml:"auto_refresh"`
}

// ConfigPath returns the default config file path
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ruslan-cli", "config.yaml")
}

// Load reads the configuration file
func Load() (*Config, error) {
	configPath := ConfigPath()

	// Create default config if doesn't exist
if _, err := os.Stat(configPath); os.IsNotExist(err) {
cfg := DefaultConfig()
if err := cfg.Save(); err != nil {
return nil, fmt.Errorf("failed to create default config: %w", err)
}
return cfg, nil
}

data, err := os.ReadFile(configPath)
if err != nil {
return nil, fmt.Errorf("failed to read config: %w", err)
}

var cfg Config
if err := yaml.Unmarshal(data, &cfg); err != nil {
return nil, fmt.Errorf("failed to parse config: %w", err)
}

return &cfg, nil
}

// Save writes the configuration file (method on Config)
func (c *Config) Save() error {
configPath := ConfigPath()
configDir := filepath.Dir(configPath)

// Create config directory if doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()

	return &Config{
		CurrentEnvironment: "dev",
	Environments: map[string]*Environment{
		"dev": {
			Name:        "Development",
			ProjectID:   "homework-475918",
			Region:      "us-central1",
			ClusterName: "dev-gke-cluster",
			Namespace:   "vault",
			ServiceName: "vault",
			VaultAddr:   "https://vault-dev.dautov.dev",
			VaultPort:   "443",
			UseNipIO:    false,
		},
		"prod": {
			Name:        "Production",
			ProjectID:   "homework-475918",
			Region:      "us-central1",
			ClusterName: "prod-gke-cluster",
			Namespace:   "vault",
			ServiceName: "vault",
			VaultAddr:   "https://vault.dautov.dev",
			VaultPort:   "443",
			UseNipIO:    false,
		},
		},
		TokenFile:    filepath.Join(home, ".ruslan-cli", "tokens"),
		CacheDir:     filepath.Join(home, ".ruslan-cli", "cache"),
		OutputFormat: "table",
		AutoRefresh:  true,
	}
}
