package config

import "strings"

type RPCServer struct {
	Namespace  string     `yaml:"namespace" env:"RPC_NAMESPACE,overwrite"`
	Name       string     `yaml:"name" env:"RPC_NAME,overwrite"`
	Version    int        `yaml:"version" env:"RPC_VERSION,overwrite"`
	Address    string     `yaml:"address" env:"RPC_ADDRESS,overwrite"`
	Resilience Resilience `yaml:"resilience"`
}

func (rps *RPCServer) Validate() error {
	rps.Namespace = strings.TrimSpace(rps.Namespace)
	rps.Name = strings.TrimSpace(rps.Name)
	rps.Address = strings.TrimSpace(rps.Address)

	if rps.Namespace == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Namespace",
			Reason:    "Should not be empty",
		}
	}

	if rps.Name == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Name",
			Reason:    "Should not be empty",
		}
	}

	if rps.Address == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}
