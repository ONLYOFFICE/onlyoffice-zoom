package shared

import (
	"strings"
)

type ZoomConfig struct {
	ClientID      string `yaml:"client_id" env:"ZOOM_CLIENT_ID,overwrite"`
	ClientSecret  string `yaml:"client_secret" env:"ZOOM_CLIENT_SECRET,overwrite"`
	WebhookSecret string `yaml:"webhook_secret" env:"ZOOM_WEBHOOK_SECRET,overwrite"`
	RedirectURI   string `yaml:"redirect_uri" env:"ZOOM_REDIRECT_URI,overwrite"`
}

func (zc *ZoomConfig) Validate() error {
	zc.ClientID = strings.TrimSpace(zc.ClientID)
	zc.ClientSecret = strings.TrimSpace(zc.ClientSecret)

	if zc.ClientID == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "ClientID",
			Reason:    "Should not be empty",
		}
	}

	if zc.ClientSecret == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "ClientSecret",
			Reason:    "Should not be empty",
		}
	}

	return nil
}

type OnlyofficeConfig struct {
	DocSecret   string `yaml:"doc_secret" env:"ONLYOFFICE_DOC_SECRET,overwrite"`
	CallbackURL string `yaml:"callback_url" env:"ONLYOFFICE_CALLBACK_URL,overwrite"`
}

func (oc *OnlyofficeConfig) Validate() error {
	oc.DocSecret = strings.TrimSpace(oc.DocSecret)

	if oc.DocSecret == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "DocSecret",
			Reason:    "Should not be empty",
		}
	}

	return nil
}

type CallbackConfig struct {
	// TODO: uint
	MaxSize       int64 `yaml:"max_size" env:"CALLBACK_MAX_SIZE,overwrite"`
	UploadTimeout int   `yaml:"upload_timeout" env:"CALLBACK_UPLOAD_TIMEOUT,overwrite"`
}

func (c *CallbackConfig) Validate() error {
	if c.MaxSize <= 0 {
		return &InvalidConfigurationParameterError{
			Parameter: "MaxSize",
			Reason:    "Must be positive",
		}
	}

	if c.UploadTimeout <= 0 {
		return &InvalidConfigurationParameterError{
			Parameter: "UploadTimeout",
			Reason:    "Must be positive",
		}
	}

	return nil
}

type RedisConfig struct {
	RedisAddresses []string `yaml:"redis_addresses" env:"REDIS_ADDRESSES,overwrite"`
	RedisUsername  string   `yaml:"redis_username" env:"REDIS_USERNAME,overwrite"`
	RedisPassword  string   `yaml:"redis_password" env:"REDIS_PASSWORD,overwrite"`
	RedisDatabase  int      `yaml:"redis_database" env:"REDIS_DATABASE,overwrite"`
	BufferSize     int      `yaml:"redis_buffer_size" env:"REDIS_BUFFER_SIZE,overwrite"`
}

func (rc *RedisConfig) Validate() error {
	if len(rc.RedisAddresses) == 0 {
		return &InvalidConfigurationParameterError{
			Parameter: "RedisAddress",
			Reason:    "Should not be empty",
		}
	}

	return nil
}
