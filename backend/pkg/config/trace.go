package config

type Tracer struct {
	Enable        bool    `yaml:"enable" env:"TRACER_ENABLE,overwrite"`
	Address       string  `yaml:"address" env:"TRACER_ADDRESS,overwrite"`
	TracerType    int     `yaml:"type" env:"TRACER_TYPE,overwrite"`
	FractionRatio float64 `yaml:"fraction" env:"TRACER_FRACTION_RATIO,overwrite"`
}
