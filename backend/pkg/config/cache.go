package config

type Cache struct {
	Type int `yaml:"type" env:"CACHE_TYPE,overwrite"`
	Size int `yaml:"size" env:"CACHE_SIZE,overwrite"`
}

func (b *Cache) Validate() error {
	return nil
}
