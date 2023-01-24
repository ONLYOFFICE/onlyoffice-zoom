package config

import "time"

type Registry struct {
	Addresses    []string      `yaml:"addresses" env:"REGISTRY_ADDRESSES,overwrite"`
	CacheTTL     time.Duration `yaml:"cache_duration" env:"REGISTRY_CACHE_DURATION,overwrite"`
	RegistryType int           `yaml:"type" env:"REGISTRY_TYPE,overwrite"`
}

func (r *Registry) Validate() error {
	return nil
}
