// Package providers implements configuration providers for the config service.
// The EnvProvider loads configuration from environment variables, supporting both
// explicit mappings (EnvName) and automatic dot-to-snake_case conversion.
package providers

import (
	"os"
	"strings"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

type EnvProvider struct{}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func (p *EnvProvider) Name() string {
	return "env"
}

// Load implements the Provider interface
func (p *EnvProvider) Load(cfg *models.Config) error {
	for i, item := range cfg.Items {
		// Priority 1: Specific EnvName set
		if item.EnvName != "" {
			envValue := os.Getenv(item.EnvName)
			if envValue != "" {
				cfg.Items[i].Value = envValue
				continue
			}
		}

		// Priority 2: Auto-mapped Env Var (UPPER_CASE with underscores)
		autoEnvName := strings.ToUpper(strings.ReplaceAll(item.Key, ".", "_"))
		envValue := os.Getenv(autoEnvName)
		if envValue != "" {
			cfg.Items[i].Value = envValue
		}
	}
	return nil
}
