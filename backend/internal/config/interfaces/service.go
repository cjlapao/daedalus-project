package interfaces

import (
	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

// ConfigService defines the interface for the configuration service
type ConfigService interface {
	Load() error
	Get() *models.Config
	// GetSection retrieves a strongly typed section from the configuration
	// matching the given key and binding it to the generic type T
	GetSection(key string, target interface{}) error
	// LoadFromFile loads a configuration from a specific file and returns the config object
	// unrelated to the main service configuration
	LoadFromFile(path string) (*models.Config, error)
}

// Provider defines the interface for configuration providers
type Provider interface {
	Name() string
	Load(cfg *models.Config) error
}
