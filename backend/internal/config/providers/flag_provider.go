// Package providers implements configuration providers for the config service.
// The FlagProvider loads configuration from command-line flags.
// It maps specific flags defined in ConfigItems to their corresponding values.
package providers

import (
	"flag"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

type FlagProvider struct{}

func NewFlagProvider() *FlagProvider {
	return &FlagProvider{}
}

func (p *FlagProvider) Name() string {
	return "flag"
}

func (p *FlagProvider) Load(cfg *models.Config) error {
	for _, item := range cfg.Items {
		if item.FlagName != "" {
			// reading the flag value
			flagValue := flag.Lookup(item.FlagName)
			if flagValue != nil {
				item.Value = flagValue.Value.String()
			}
		}
	}

	return nil
}
