package config

type Cache struct {
	Type     int    `yaml:"type" env:"CACHE_TYPE,overwrite"`
	Size     int    `yaml:"size" env:"CACHE_SIZE,overwrite"`
	Address  string `yaml:"address" env:"CACHE_ADDRESS,overwrite"`
	Password string `yaml:"password" env:"CACHE_PASSWORD,overwrite"`
	DB       int    `yaml:"db" env:"CACHE_DB,overwrite"`
}

func (b *Cache) Validate() error {
	return nil
}
