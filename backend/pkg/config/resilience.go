package config

type RateLimiter struct {
	Limit   uint64 `yaml:"limit" env:"RATE_LIMIT,overwrite"`
	IPLimit uint64 `yaml:"iplimit" env:"RATE_LIMIT_IP,overwrite"`
}
