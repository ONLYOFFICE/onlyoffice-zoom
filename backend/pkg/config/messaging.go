package config

type Broker struct {
	Type   int      `yaml:"type" env:"BROKER_TYPE,overwrite"`
	Addrs  []string `yaml:"addresses" env:"BROKER_ADDRESSES,overwrite"`
	Secure bool     `yaml:"secure" env:"BROKER_SECURE,overwrite"`
}

func (b *Broker) Validate() error {
	if len(b.Addrs) == 0 {
		return &InvalidConfigurationParameterError{
			Parameter: "Addrs",
			Reason:    "Invalid number of addresses",
		}
	}

	return nil
}
