package config

// Worker configuration
type WorkerConfig struct {
	MaxConcurrency int      `yaml:"max_concurrency" env:"WORKER_MAX_CONCURRENCY,overwrite"`
	RedisAddresses []string `yaml:"addresses" env:"WORKER_ADDRESS,overwrite"`
	RedisUsername  string   `yaml:"username" env:"WORKER_USERNAME,overwrite"`
	RedisPassword  string   `yaml:"password" env:"WORKER_PASSWORD,overwrite"`
	RedisDatabase  int      `yaml:"database" env:"WORKER_DATABASE,overwrite"`
}

func (wc *WorkerConfig) Validate() error {
	if len(wc.RedisAddresses) < 1 {
		return &InvalidConfigurationParameterError{
			Parameter: "Worker address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}
