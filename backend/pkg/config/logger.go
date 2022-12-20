package config

type Logger struct {
	Name    string     `yaml:"name" env:"LOGGER_NAME,overwrite"`
	Level   int        `yaml:"level" env:"LOGGER_LEVEL,overwrite"`
	Pretty  bool       `yaml:"pretty" env:"LOGGER_PRETTY,overwrite"`
	Color   bool       `yaml:"color" env:"LOGGER_COLOR,overwrite"`
	File    FileLog    `yaml:"file"`
	Elastic ElasticLog `yaml:"elastic"`
}

type ElasticLog struct {
	Address            string `yaml:"address" env:"ELASTIC_ADDRESS,overwrite"`
	Index              string `yaml:"index" env:"ELASTIC_INDEX,overwrite"`
	Level              int    `yaml:"level" env:"ELASTIC_LEVEL,overwrite"`
	Bulk               bool   `yaml:"bulk" env:"ELASTIC_BULK,overwrite"`
	Async              bool   `yaml:"async" env:"ELASTIC_ASYNC,overwrite"`
	HealthcheckEnabled bool   `yaml:"healthcheck" env:"ELASTIC_HEALTHCHECK,overwrite"`
	BasicAuthUsername  string `yaml:"username" env:"ELASTIC_AUTH_USERNAME,overwrite"`
	BasicAuthPassword  string `yaml:"password" env:"ELASTIC_AUTH_PASSWORD,overwrite"`
	GzipEnabled        bool   `yaml:"gzip" env:"ELASTIC_GZIP_ENABLED,overwrite"`
}

type FileLog struct {
	Filename   string `yaml:"filename" env:"FILELOG_NAME,overwrite"`
	MaxSize    int    `yaml:"maxsize" env:"FILELOG_MAX_SIZE,overwrite"`
	MaxAge     int    `yaml:"maxage" env:"FILELOG_MAX_AGE,overwrite"`
	MaxBackups int    `yaml:"maxbackups" env:"FILELOG_MAX_BACKUPS,overwrite"`
	LocalTime  bool   `yaml:"localtime"`
	Compress   bool   `yaml:"compress" env:"FILELOG_COMPRESS,overwrite"`
}
