package config

import "strings"

// Worker configuration
type WorkerConfig struct {
	MaxConcurrency int    `yaml:"max_concurrency" env:"WORKER_MAX_CONCURRENCY,overwrite"`
	RedisAddress   string `yaml:"address" env:"WORKER_ADDRESS,overwrite"`
	RedisUsername  string `yaml:"username" env:"WORKER_USERNAME,overwrite"`
	RedisPassword  string `yaml:"password" env:"WORKER_PASSWORD,overwrite"`
	RedisDatabase  int    `yaml:"database" env:"WORKER_DATABASE,overwrite"`
}

func (wc *WorkerConfig) Validate() error {
	wc.RedisAddress = strings.TrimSpace(wc.RedisAddress)

	if wc.RedisAddress == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Worker address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}
