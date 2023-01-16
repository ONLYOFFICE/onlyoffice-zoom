package config

type Resilience struct {
	RateLimiter    RateLimiter `yaml:"rate_limiter"`
	CircuitBreaker Hystrix     `yaml:"circuit_breaker"`
}

type RateLimiter struct {
	Limit   uint64 `yaml:"limit" env:"RATE_LIMIT,overwrite"`
	IPLimit uint64 `yaml:"iplimit" env:"RATE_LIMIT_IP,overwrite"`
}

// Go-micro hystrix configuration
type Hystrix struct {
	// Timeout is how long to wait for command to complete, in milliseconds (default 1000)
	Timeout int `yaml:"timeout" env:"HYSTRIX_TIMEOUT,overwrite"`
	// MaxConcurrent is how many commands of the same type can run at the same time (default 10)
	MaxConcurrent int `yaml:"max_concurrent" env:"HYSTRIX_MAX_CONCURRENT,overwrite"`
	// VolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health (default 20)
	VolumeThreshold int `yaml:"volume_threshold" env:"HYSTRIX_VOLUME_THRESHOLD,overwrite"`
	// SleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery (default 5000)
	SleepWindow int `yaml:"sleep_window" env:"HYSTRIX_SLEEP_WINDOW,overwrite"`
	// ErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests (default 50)
	ErrorPercentThreshold int `yaml:"error_percent_threshold" env:"HYSTRIX_ERROR_PERCENT_THRESHOLD,overwrite"`
}
