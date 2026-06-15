// Package providers implements configuration providers for the config service.
// The FileProvider loads configuration from YAML or JSON files.
// It supports recursive flattening of nested structures into dot-notation keys
// (e.g. "server.port") for consistent access.
package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
	"gopkg.in/yaml.v3"
)

// FileProvider provides configuration from a JSON file
type FileProvider struct {
	filePath string
}

// NewFileProvider creates a new file provider
func NewFileProvider(filePath string) *FileProvider {
	return &FileProvider{
		filePath: filePath,
	}
}

func (p *FileProvider) Name() string {
	return "file"
}

// Load implements the Provider interface
func (p *FileProvider) Load(cfg *models.Config) error {
	path := p.resolveFilePath()
	if path == "" {
		// No config file found — not an error, other providers will supply values
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse based on extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yaml" || ext == ".yml" {
		var yamlMap map[string]interface{}
		if err := yaml.Unmarshal(data, &yamlMap); err != nil {
			return fmt.Errorf("failed to parse yaml config file: %w", err)
		}
		p.flattenMap("", yamlMap, cfg)
	} else {
		// Default to JSON but handle potentially new structure or flat structure
		// For backward compatibility mostly, but ideally we treat JSON same as YAML (nested or flat)
		// Try unmarshalling as simple key-value list first (legacy)
		var fileConfig []models.ConfigItem
		if err := json.Unmarshal(data, &fileConfig); err == nil {
			// Legacy format: array of ConfigItem
			cfg.Items = append(cfg.Items, fileConfig...)
			return nil
		}

		// If that fails, try map
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(data, &jsonMap); err != nil {
			return fmt.Errorf("failed to parse config file as legacy list or json map: %w", err)
		}
		p.flattenMap("", jsonMap, cfg)
	}

	return nil
}

// resolveFilePath determines the config file path using the following priority:
// 1. Explicitly provided path (constructor argument)
// 2. CONFIG_FILE_PATH environment variable
// 3. --config-path command-line flag
// 4. config.yaml or config.yml next to the executable
func (p *FileProvider) resolveFilePath() string {
	// 1. Explicit path from constructor
	if p.filePath != "" {
		return p.filePath
	}

	// 2. Environment variable
	if envPath := os.Getenv(models.ConfigFilePathEnv); envPath != "" {
		return envPath
	}

	// 3. Command-line argument (scan os.Args directly since flag.Parse
	//    hasn't run yet — this is a bootstrap-level flag)
	if argPath := findArgValue("--" + models.ConfigFilePathFlag); argPath != "" {
		return argPath
	}

	// 4. Probe for config.yaml / config.yml next to the executable
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	execDir := filepath.Dir(execPath)

	for _, name := range []string{"config.yaml", "config.yml"} {
		candidate := filepath.Join(execDir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// 5. Probe in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for _, name := range []string{"config.yaml", "config.yml"} {
		candidate := filepath.Join(cwd, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

func (p *FileProvider) flattenMap(prefix string, data map[string]interface{}, cfg *models.Config) {
	for k, v := range data {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		p.flattenValue(fullKey, v, cfg)
	}
}

func (p *FileProvider) flattenValue(key string, value interface{}, cfg *models.Config) {
	switch val := value.(type) {
	case map[string]interface{}:
		p.flattenMap(key, val, cfg)
	case map[interface{}]interface{}: // YAML often produces this
		strMap := make(map[string]interface{})
		for mk, mv := range val {
			strMap[fmt.Sprintf("%v", mk)] = mv
		}
		p.flattenMap(key, strMap, cfg)
	case []interface{}:
		for i, v := range val {
			p.flattenValue(fmt.Sprintf("%s.%d", key, i), v, cfg)
		}
	default:
		// Convert primitive to string using Set
		cfg.Set(key, fmt.Sprintf("%v", val))
	}
}

func (p *FileProvider) IsAvailable() bool {
	return p.resolveFilePath() != ""
}

// findArgValue scans os.Args for a flag in either "--flag value" or "--flag=value" form.
func findArgValue(flag string) string {
	args := os.Args[1:]
	for i, arg := range args {
		// --config-path=value
		if strings.HasPrefix(arg, flag+"=") {
			return strings.TrimPrefix(arg, flag+"=")
		}
		// --config-path value
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
