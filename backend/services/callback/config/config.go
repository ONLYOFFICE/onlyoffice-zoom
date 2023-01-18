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
	Server      config.HttpServer       `yaml:"server"`
	REPL        config.ReplServer       `yaml:"repl"`
	Onlyoffice  shared.OnlyofficeConfig `yaml:"onlyoffice"`
	Broker      config.Broker           `yaml:"broker"`
	Registry    config.Registry         `yaml:"registry"`
	Logger      config.Logger           `yaml:"logger"`
	Tracer      config.Tracer           `yaml:"tracer"`
	Worker      config.WorkerConfig     `yaml:"worker"`
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

	if err := c.Broker.Validate(); err != nil {
		return err
	}

	if err := c.Worker.Validate(); err != nil {
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
		Server: config.HttpServer{
			Namespace: "onlyoffices",
			Name:      "callback",
			Address:   ":4141",
		},
		REPL: config.ReplServer{
			Namespace: "onlyoffice",
			Name:      "callback.repl",
			Version:   0,
			Address:   ":4242",
			Debug:     true,
		},
		Context: context.Background(),
	}

	return config
}
