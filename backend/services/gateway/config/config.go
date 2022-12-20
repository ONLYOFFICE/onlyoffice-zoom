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
	Server      config.HttpServer `yaml:"server"`
	Zoom        shared.ZoomConfig `yaml:"zoom"`
	REPL        config.ReplServer `yaml:"repl"`
	Broker      config.Broker     `yaml:"broker"`
	Registry    config.Registry   `yaml:"registry"`
	Logger      config.Logger     `yaml:"logger"`
	Tracer      config.Tracer     `yaml:"tracer"`
	Context     context.Context   `yaml:"-"`
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
		Server: config.HttpServer{
			Namespace: "onlyoffices",
			Name:      "gateway",
			Address:   ":9191",
		},
		REPL: config.ReplServer{
			Namespace: "onlyoffice",
			Name:      "gateway",
			Version:   0,
			Address:   ":9992",
			Debug:     true,
		},
		Zoom: shared.ZoomConfig{
			WebhookSecret: "secret",
			RedirectURI:   "https://onlyoffice.com",
		},
		Context: context.Background(),
	}

	return config
}
