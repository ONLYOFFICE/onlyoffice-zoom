package config

import "strings"

type HttpServer struct {
	Namespace   string      `yaml:"namespace" env:"HTTP_NAMESPACE,overwrite"`
	Name        string      `yaml:"name" env:"HTTP_NAME,overwrite"`
	Version     int         `yaml:"version" env:"HTTP_VERSION,overwrite"`
	Address     string      `yaml:"address" env:"HTTP_ADDRESS,overwrite"`
	RateLimiter RateLimiter `yaml:"rate_limiter"`
	CORS        CORS        `yaml:"cors"`
}

func (hs *HttpServer) Validate() error {
	hs.Namespace = strings.TrimSpace(hs.Namespace)
	hs.Name = strings.TrimSpace(hs.Name)
	hs.Address = strings.TrimSpace(hs.Address)

	if hs.Namespace == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Namespace",
			Reason:    "Should not be empty",
		}
	}

	if hs.Name == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Name",
			Reason:    "Should not be empty",
		}
	}

	if hs.Address == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}

type CORS struct {
	AllowedOrigins     []string `yaml:"origins" env:"ALLOWED_ORIGINS,overwrite"`
	AllowedMethods     []string `yaml:"methods" env:"ALLOWED_METHODS,overwrite"`
	AllowedHeaders     []string `yaml:"headers" env:"ALLOWED_HEADERS,overwrite"`
	AllowedCredentials bool     `yaml:"credentials" env:"ALLOW_CREDENTIALS,overwrite"`
}
