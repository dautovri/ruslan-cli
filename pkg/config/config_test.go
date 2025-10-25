package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.CurrentEnvironment != "dev" {
		t.Errorf("expected default environment 'dev', got %s", cfg.CurrentEnvironment)
	}

	if len(cfg.Environments) == 0 {
		t.Error("expected environments to be populated")
	}

	// Check dev environment exists
	if _, ok := cfg.Environments["dev"]; !ok {
		t.Error("expected 'dev' environment to exist")
	}

	// Check prod environment exists
	if _, ok := cfg.Environments["prod"]; !ok {
		t.Error("expected 'prod' environment to exist")
	}
}

func TestConfigPath(t *testing.T) {
	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath failed: %v", err)
	}

	if path == "" {
		t.Error("configPath returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("configPath should return absolute path, got: %s", path)
	}

	if filepath.Base(filepath.Dir(path)) != ".ruslan-cli" {
		t.Errorf("config should be in .ruslan-cli directory, got: %s", path)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	cfg := DefaultConfig()
	cfg.CurrentEnvironment = "prod"
	cfg.Environments["dev"].Token = "test-token"

	// Override config path for testing
	oldConfigPath := configPath
	configPath = func() (string, error) {
		return configFile, nil
	}
	defer func() { configPath = oldConfigPath }()

	// Save config
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded config matches saved config
	if loadedCfg.CurrentEnvironment != cfg.CurrentEnvironment {
		t.Errorf("CurrentEnvironment mismatch: got %s, want %s",
			loadedCfg.CurrentEnvironment, cfg.CurrentEnvironment)
	}

	if loadedCfg.Environments["dev"].Token != cfg.Environments["dev"].Token {
		t.Errorf("Token mismatch: got %s, want %s",
			loadedCfg.Environments["dev"].Token, cfg.Environments["dev"].Token)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "nonexistent.yaml")

	// Override config path
	oldConfigPath := configPath
	configPath = func() (string, error) {
		return configFile, nil
	}
	defer func() { configPath = oldConfigPath }()

	// Load should create default config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should create default config, got error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load returned nil config")
	}

	// Verify it's the default config
	if cfg.CurrentEnvironment != "dev" {
		t.Errorf("expected default environment 'dev', got %s", cfg.CurrentEnvironment)
	}
}

func TestEnvironmentValidation(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name    string
		envName string
		wantErr bool
	}{
		{
			name:    "valid dev environment",
			envName: "dev",
			wantErr: false,
		},
		{
			name:    "valid prod environment",
			envName: "prod",
			wantErr: false,
		},
		{
			name:    "invalid environment",
			envName: "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := cfg.Environments[tt.envName]
			if tt.wantErr && exists {
				t.Error("expected environment not to exist")
			}
			if !tt.wantErr && !exists {
				t.Error("expected environment to exist")
			}
		})
	}
}
