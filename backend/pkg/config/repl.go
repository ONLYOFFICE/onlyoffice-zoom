package config

import "strings"

type ReplServer struct {
	Namespace string `yaml:"namespace" env:"REPL_NAMESPACE,overwrite"`
	Name      string `yaml:"name" env:"REPL_NAME,overwrite"`
	Version   int    `yaml:"version" env:"REPL_VERSION,overwrite"`
	Address   string `yaml:"address" env:"REPL_ADDRESS,overwrite"`
	Debug     bool   `yaml:"debug" env:"REPL_DEBUG,overwrite"`
}

func (rs *ReplServer) Validate() error {
	rs.Namespace = strings.TrimSpace(rs.Namespace)
	rs.Name = strings.TrimSpace(rs.Name)
	rs.Address = strings.TrimSpace(rs.Address)

	if rs.Namespace == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Namespace",
			Reason:    "Should not be empty",
		}
	}

	if rs.Name == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Name",
			Reason:    "Should not be empty",
		}
	}

	if rs.Address == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}
