package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Env  string `envconfig:"PIGLET_ENV" default:"prod"`
	GRPC GRPCConfig
}

type GRPCConfig struct {
	Server string `evconfig:"PIGLET_BILLS_SERVER" default:"localhost"`
	Port   int    `envconfig:"PIGLET_BILLS_PORT" default:"8080"`
}

// InitConfig reads config variables from env and init *Config value
func MustLoadConfig() *Config {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		panic("piglet-bills: failed to read config: " + err.Error())
	}

	return cfg
}
