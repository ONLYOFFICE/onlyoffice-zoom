package config

import (
	"context"
	"os"
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared"
)

type Config struct {
	Environment string `yaml:"environment" env:"ENVIRONMENT,overwrite"`
	Machine     string
	Server      config.RPCServer        `yaml:"server"`
	REPL        config.ReplServer       `yaml:"repl"`
	Zoom        shared.ZoomConfig       `yaml:"zoom"`
	Onlyoffice  shared.OnlyofficeConfig `yaml:"onlyoffice"`
	Redis       shared.RedisConfig      `yaml:"redis"`
	Broker      config.Broker           `yaml:"broker"`
	Registry    config.Registry         `yaml:"registry"`
	Logger      config.Logger           `yaml:"logger"`
	Tracer      config.Tracer           `yaml:"tracer"`
	Cache       config.Cache            `yaml:"cache"`
	Context     context.Context         `yaml:"-"`
}

func (c *Config) Validate() error {
	c.Environment = strings.TrimSpace(c.Environment)
	c.Machine = strings.TrimSpace(c.Machine)

	if c.Environment == "" {
		return &shared.InvalidConfigurationParameterError{
			Parameter: "Environemnt",
			Reason:    "Should not be empty",
		}
	}

	if c.Machine == "" {
		return &shared.InvalidConfigurationParameterError{
			Parameter: "Machine",
			Reason:    "Should not be empty",
		}
	}

	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.REPL.Validate(); err != nil {
		return err
	}

	if err := c.Onlyoffice.Validate(); err != nil {
		return err
	}

	if err := c.Zoom.Validate(); err != nil {
		return err
	}

	if err := c.Broker.Validate(); err != nil {
		return err
	}

	if err := c.Registry.Validate(); err != nil {
		return err
	}

	return nil
}

func BuildConfig() *Config {
	host, _ := os.Hostname()

	if host == "" {
		host = "anonymous-machine"
	}

	config := &Config{
		Environment: "production",
		Machine:     host,
		Server: config.RPCServer{
			Namespace: "onlyoffice",
			Name:      "builder",
			Address:   ":6060",
		},
		REPL: config.ReplServer{
			Namespace: "onlyoffice",
			Name:      "builder.repl",
			Version:   0,
			Address:   ":7979",
			Debug:     true,
		},
		Context: context.Background(),
	}

	return config
}
