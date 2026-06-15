package config

import (
	"fmt"

	"github.com/cjlapao/daedalus/backend/internal/config/interfaces"
	"github.com/cjlapao/daedalus/backend/internal/config/models"
	"github.com/cjlapao/daedalus/backend/internal/config/service"
)

// Re-export common types
type (
	Config        = models.Config
	ConfigItem    = models.ConfigItem
	ConfigService = service.ConfigService
	Provider      = interfaces.Provider
)

// Initialize initializes the configuration service
func Initialize(providers ...interfaces.Provider) (*service.ConfigService, error) {
	return service.Initialize(providers...)
}

// GetInstance returns the service instance
func GetInstance() *service.ConfigService {
	return service.GetInstance()
}

// Reset clears the singleton instance (for testing)
func Reset() {
	service.Reset()
}

// GetSection retrieves a strongly typed section from the configuration
// matching the given key and binding it to the generic type T.
// This is a helper wrapper around the service's GetSection method.
func GetSection[T any](key string) (T, error) {
	var result T
	svc := service.GetInstance()
	if svc == nil {
		return result, fmt.Errorf("config service not initialized")
	}

	if err := svc.GetSection(key, &result); err != nil {
		return result, err
	}

	return result, nil
}

// LoadFromFile loads a configuration from a specific file
func LoadFromFile(path string) (*models.Config, error) {
	svc := service.GetInstance()
	if svc == nil {
		return service.New().LoadFromFile(path)
	}
	return svc.LoadFromFile(path)
}
