package config

import "time"

type Registry struct {
	Addresses    []string      `yaml:"addresses" env:"REGISTRY_ADDRESSES,overwrite"`
	Secure       bool          `yaml:"secure" env:"REGISTRY_SECURE,overwrite"`
	CacheTTL     time.Duration `yaml:"cache_duration" env:"REGISTRY_CACHE_DURATION,overwrite"`
	RegistryType int           `yaml:"type" env:"REGISTRY_TYPE,overwrite"`
}

func (r *Registry) Validate() error {
	if len(r.Addresses) == 0 {
		return &InvalidConfigurationParameterError{
			Parameter: "Addresses",
			Reason:    "Invalid number of addresses",
		}
	}

	return nil
}
