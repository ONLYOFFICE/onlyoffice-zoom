package config

type Broker struct {
	Addrs          []string `yaml:"addresses" env:"BROKER_ADDRESSES,overwrite"`
	Type           int      `yaml:"type" env:"BROKER_TYPE,overwrite"`
	DisableAutoAck bool     `yaml:"disable_auto_ack" env:"BROKER_DISABLE_AUTO_ACK,overwrite"`
	Durable        bool     `yaml:"durable" env:"BROKER_DURABLE,overwrite"`
	AckOnSuccess   bool     `yaml:"ack_on_success" env:"BROKER_ACK_ON_SUCCESS,overwrite"`
	RequeueOnError bool     `yaml:"requeue_on_error" env:"BROKER_REQUEUE_ON_ERROR,overwrite"`
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
