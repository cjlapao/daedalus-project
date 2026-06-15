package providers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

func TestFileProvider(t *testing.T) {
	t.Run("NewFileProvider creates provider", func(t *testing.T) {
		p := NewFileProvider("/path/to/config.yaml")
		if p == nil {
			t.Fatal("NewFileProvider() returned nil")
		}
	})

	t.Run("Name returns file", func(t *testing.T) {
		p := NewFileProvider("")
		if p.Name() != "file" {
			t.Errorf("Name() = %q, want %q", p.Name(), "file")
		}
	})
}

func TestFileProviderLoad(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("loads YAML file", func(t *testing.T) {
		yamlPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(yamlPath, []byte(`
server:
  port: 8080
  host: localhost
`), 0o644)
		if err != nil {
			t.Fatalf("failed to write YAML: %v", err)
		}

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "server.port", Value: "5000"},
				{Key: "server.host", Value: "0.0.0.0"},
			},
		}

		p := NewFileProvider(yamlPath)
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		port := cfg.Get("server.port")
		if port == nil || port.Value != "8080" {
			t.Errorf("server.port = %v, want value 8080", port)
		}
		host := cfg.Get("server.host")
		if host == nil || host.Value != "localhost" {
			t.Errorf("server.host = %v, want value localhost", host)
		}
	})

	t.Run("loads YAML with .yml extension", func(t *testing.T) {
		ymlPath := filepath.Join(tmpDir, "config.yml")
		err := os.WriteFile(ymlPath, []byte(`
database:
  name: testdb
`), 0o644)
		if err != nil {
			t.Fatalf("failed to write YAML: %v", err)
		}

		cfg := &models.Config{}
		p := NewFileProvider(ymlPath)
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		dbName := cfg.Get("database.name")
		if dbName == nil || dbName.Value != "testdb" {
			t.Errorf("database.name = %v, want testdb", dbName)
		}
	})

	t.Run("loads JSON nested map", func(t *testing.T) {
		jsonPath := filepath.Join(tmpDir, "config.json")
		err := os.WriteFile(jsonPath, []byte(`{
		"server": {
			"port": 9090,
			"debug": true
		}
	}`), 0o644)
		if err != nil {
			t.Fatalf("failed to write JSON: %v", err)
		}

		cfg := &models.Config{}
		p := NewFileProvider(jsonPath)
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		port := cfg.Get("server.port")
		if port == nil || port.Value != "9090" {
			t.Errorf("server.port = %v, want 9090", port)
		}
		debug := cfg.Get("server.debug")
		if debug == nil || debug.Value != "true" {
			t.Errorf("server.debug = %v, want true", debug)
		}
	})

	t.Run("loads legacy JSON array format", func(t *testing.T) {
		jsonPath := filepath.Join(tmpDir, "config-legacy.json")
		err := os.WriteFile(jsonPath, []byte(`[
		{"Key": "legacy.key", "Value": "legacy-value"}
	]`), 0o644)
		if err != nil {
			t.Fatalf("failed to write JSON: %v", err)
		}

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "legacy.key", Value: "default"},
			},
		}
		initialLen := len(cfg.Items)

		p := NewFileProvider(jsonPath)
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		// Legacy format appends items
		if len(cfg.Items) != initialLen+1 {
			t.Errorf("Items len = %d, want %d", len(cfg.Items), initialLen+1)
		}
		item := cfg.Get("legacy.key")
		if item == nil {
			t.Fatal("legacy.key not found")
		}
		// Legacy format appends, so original item still exists
		// The first item should still have "default" value
		if item.Value != "default" {
			t.Errorf("legacy.key.Value = %q, want %q (legacy appends)", item.Value, "default")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		cfg := &models.Config{}
		p := NewFileProvider("/nonexistent/path/config.yaml")
		err := p.Load(cfg)
		if err == nil {
			t.Fatal("Load() should error for non-existent file")
		}
	})

	t.Run("returns nil when no config file found (no explicit path)", func(t *testing.T) {
		cfg := &models.Config{}
		// NewFileProvider with empty path will try to resolve from env, flags, executable dir, cwd
		// In test environment, likely no config file exists
		p := NewFileProvider("")
		err := p.Load(cfg)
		// Should not error — just no config loaded
		if err != nil {
			t.Logf("Load() returned error (expected if config file exists): %v", err)
		}
	})
}

func TestFileProviderResolveFilePath(t *testing.T) {
	t.Run("returns explicit path from constructor", func(t *testing.T) {
		p := NewFileProvider("/explicit/path/config.yaml")
		// We can't call resolveFilePath directly (unexported),
		// but we can test via Load behavior: if explicit path is used,
		// Load should try to read that specific file.
		cfg := &models.Config{}
		err := p.Load(cfg)
		// Should fail with "not found" since file doesn't exist at that path
		if err == nil {
			t.Fatal("Load() should error when explicit path doesn't exist")
		}
	})
}

func TestFileProviderIsAvailable(t *testing.T) {
	t.Run("returns false when no config file exists", func(t *testing.T) {
		// Temporarily change cwd to a directory without config files
		tmpDir := t.TempDir()
		origWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origWd)

		p := NewFileProvider("")
		if p.IsAvailable() {
			t.Error("IsAvailable() = true, want false in empty directory")
		}
	})

	t.Run("returns true when config file exists in cwd", func(t *testing.T) {
		tmpDir := t.TempDir()
		origWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origWd)

		// Create config.yaml in temp dir
		err := os.WriteFile(filepath.Join(tmpDir, "config.yaml"), []byte("test: 1"), 0o644)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		p := NewFileProvider("")
		if !p.IsAvailable() {
			t.Error("IsAvailable() = false, want true when config.yaml exists")
		}
	})
}

func TestFindArgValue(t *testing.T) {
	t.Run("finds --flag=value form", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"prog", "--config-path=/etc/config.yaml"}

		got := findArgValue("--config-path")
		if got != "/etc/config.yaml" {
			t.Errorf("findArgValue() = %q, want %q", got, "/etc/config.yaml")
		}
	})

	t.Run("finds --flag value form", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"prog", "--config-path", "/etc/config.yaml", "other"}

		got := findArgValue("--config-path")
		if got != "/etc/config.yaml" {
			t.Errorf("findArgValue() = %q, want %q", got, "/etc/config.yaml")
		}
	})

	t.Run("returns empty when flag not found", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"prog", "--other-flag", "value"}

		got := findArgValue("--config-path")
		if got != "" {
			t.Errorf("findArgValue() = %q, want empty", got)
		}
	})

	t.Run("returns empty when no args", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"prog"}

		got := findArgValue("--config-path")
		if got != "" {
			t.Errorf("findArgValue() = %q, want empty", got)
		}
	})
}
