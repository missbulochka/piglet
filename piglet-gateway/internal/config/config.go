package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GRPC GRPCConfig
}

type GRPCConfig struct {
	GatewayServer string `envconfig:"PIGLET_GATEWAY_SERVER" default:"0.0.0.0"`
	GatewayPort   string `envconfig:"PIGLET_GATEWAY_PORT" default:"8083"`
	BillsServer   string `envconfig:"PIGLET_BILLS_SERVER" default:"piglet-bills"`
	BillsPort     string `envconfig:"PIGLET_BILLS_PORT" default:"8080"`
	TransServer   string `envconfig:"PIGLET_TRANSACTIONS_SERVER" default:"piglet-transactions"`
	TransPort     string `envconfig:"PIGLET_TRANSACTIONS_PORT" default:"8081"`
}

// InitConfig reads config variables from env and init *Config value
func MustLoadConfig() *Config {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		panic("piglet-bills: failed to read config: " + err.Error())
	}

	return cfg
}
