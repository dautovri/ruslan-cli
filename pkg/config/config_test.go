package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.CurrentEnvironment != "dev" {
		t.Errorf("Expected default environment=dev, got: %s", cfg.CurrentEnvironment)
	}

	if _, ok := cfg.Environments["dev"]; !ok {
		t.Error("Expected 'dev' environment in default config")
	}
	if _, ok := cfg.Environments["prod"]; !ok {
		t.Error("Expected 'prod' environment in default config")
	}
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath should return absolute path, got: %s", path)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	testConfigPath := filepath.Join(tempDir, "config.yaml")

	cfg := DefaultConfig()
	cfg.CurrentEnvironment = "prod"

	configDir := filepath.Dir(testConfigPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(testConfigPath, data, 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	loadedData, err := os.ReadFile(testConfigPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loaded Config
	if err := yaml.Unmarshal(loadedData, &loaded); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if loaded.CurrentEnvironment != "prod" {
		t.Errorf("Expected current_environment=prod, got: %s", loaded.CurrentEnvironment)
	}
}

func TestEnvironmentValidation(t *testing.T) {
	cfg := DefaultConfig()

	if _, ok := cfg.Environments["dev"]; !ok {
		t.Error("Expected 'dev' environment to exist")
	}

	if _, ok := cfg.Environments["prod"]; !ok {
		t.Error("Expected 'prod' environment to exist")
	}

	if _, ok := cfg.Environments["nonexistent"]; ok {
		t.Error("Expected 'nonexistent' environment to not exist")
	}
}
