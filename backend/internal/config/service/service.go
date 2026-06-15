package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cjlapao/daedalus/backend/internal/config/interfaces"
	"github.com/cjlapao/daedalus/backend/internal/config/models"
	"github.com/cjlapao/daedalus/backend/internal/config/providers"
)

var (
	instance *ConfigService
	once     sync.Once
)

// ConfigService represents the configuration service
type ConfigService struct {
	config      *models.Config
	storagePath string
	providers   []interfaces.Provider
	isLoaded    bool
	mutex       sync.RWMutex
}

// GetInstance returns the singleton instance of the config service
func GetInstance() *ConfigService {
	if instance == nil {
		return nil
	}
	return instance
}

func New() *ConfigService {
	return &ConfigService{
		config:    models.DefaultConfig(),
		providers: make([]interfaces.Provider, 0),
	}
}

// Initialize initializes the config service singleton.
// If no providers are specified, default providers (file, env, flag) are added.
func Initialize(initialProviders ...interfaces.Provider) (*ConfigService, error) {
	var initErr error
	once.Do(func() {
		instance = New()

		if len(initialProviders) == 0 {
			// Add default providers: file (default path), env, flag
			instance.providers = append(instance.providers,
				providers.NewFileProvider(""),
				providers.NewEnvProvider(),
				providers.NewFlagProvider(),
			)
		} else {
			for _, p := range initialProviders {
				instance.AddProvider(p)
			}
		}

		if err := instance.Load(); err != nil {
			initErr = err
		}
	})
	return instance, initErr
}

// Reset clears the singleton instance (for testing)
func Reset() {
	instance = nil
	once = sync.Once{}
}

// --- System Service Interface Implementation ---

func (s *ConfigService) Name() string {
	return "config"
}

func (s *ConfigService) Init(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Return early if already loaded (idempotent) - prevents double-loading
	// when system.Start() calls Init after Initialize() already loaded config
	if s.isLoaded {
		return nil
	}

	if len(s.providers) == 0 {
		s.providers = append(s.providers,
			providers.NewFileProvider(""),
			providers.NewEnvProvider(),
			providers.NewFlagProvider(),
		)
	}

	return s.loadInternal()
}

func (s *ConfigService) Health(ctx context.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if !s.isLoaded {
		return fmt.Errorf("config service not loaded")
	}
	return nil
}

func (s *ConfigService) IsEnabled() bool {
	return true
}

func (s *ConfigService) Dependencies() []string {
	return []string{} // Config usually has no dependencies
}

// --- ConfigService Interface Implementation ---

func (s *ConfigService) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.loadInternal()
}

func (s *ConfigService) loadInternal() error {
	// Reset config to defaults before loading
	s.config = models.DefaultConfig()

	for _, provider := range s.providers {
		if err := provider.Load(s.config); err != nil {
			return fmt.Errorf("failed to load config from provider %s: %w", provider.Name(), err)
		}
	}

	// Post-load setup (e.g. storage path)
	if err := s.setupStoragePath(); err != nil {
		return err
	}

	s.isLoaded = true
	return nil
}

func (s *ConfigService) Get() *models.Config {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.config
}

func (s *ConfigService) GetSection(key string, target interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.config == nil {
		return fmt.Errorf("config not loaded")
	}

	return s.config.Bind(context.Background(), key, target)
}

func (s *ConfigService) LoadFromFile(path string) (*models.Config, error) {
	// Create a temporary config object
	cfg := models.DefaultConfig()

	provider := providers.NewFileProvider(path)
	if err := provider.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// --- Additional Helper Methods ---

func (s *ConfigService) AddProvider(provider interfaces.Provider) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.providers = append(s.providers, provider)
	s.isLoaded = false // Require reload
}

func (s *ConfigService) setupStoragePath() error {
	// Current logic from previous service.go
	// Get user folder
	userFolder, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Check if storage path is set in config
	storagePath := s.config.GetString(context.Background(), models.DatabaseStoragePathKey, "")
	if storagePath == "" {
		storagePath = filepath.Join(userFolder, models.SystemStoragePath)
	}

	// creating the folder if it doesn't exist
	if err := os.MkdirAll(storagePath, 0o755); err != nil {
		return fmt.Errorf("failed to create locally storage directory: %w", err)
	}
	s.storagePath = storagePath
	return nil
}
