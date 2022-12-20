package config

type RedisWorkerConfig struct {
	MaxActive      int    `yaml:"max_active" env:"WORKER_MAX_ACTIVE,overwrite"`
	MaxIdle        int    `yaml:"max_size" env:"WORKER_MAX_SIZE,overwrite"`
	MaxConcurrency uint   `yaml:"max_concurrency" env:"WORKER_MAX_CONCURRENCY,overwrite"`
	RedisNamespace string `yaml:"namespace" env:"WORKER_NAMESPACE,overwrite"`
	RedisAddress   string `yaml:"address" env:"WORKER_ADDRESS,overwrite"`
	RedisUsername  string `yaml:"username" env:"WORKER_USERNAME,overwrite"`
	RedisPassword  string `yaml:"password" env:"WORKER_PASSWORD,overwrite"`
	RedisDatabase  int    `yaml:"database" env:"WORKER_DATABASE,overwrite"`
}
