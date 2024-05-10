package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Env       string `envconfig:"PIGLET_ENV" default:"prod"`
	GRPC      GRPCConfig
	DB        DataBaseConfig
	RMQConfig RabbitMQConfig
}

type GRPCConfig struct {
	GRPCServer string `envconfig:"PIGLET_BILLS_SERVER" default:"0.0.0.0"`
	GRPCPort   string `envconfig:"PIGLET_BILLS_PORT" default:"8080"`
}

type DataBaseConfig struct {
	MigrationPath string `envconfig:"PIGLET_BILLS_MIGRATION_PATH" default:"./migration"`
	UserName      string `envconfig:"PIGLET_BILLS_USER" default:"postgres"`
	Password      string `envconfig:"PIGLET_BILLS_PASSWORD" default:"pass1234"`
	DBHost        string `envconfig:"PIGLET_BILLS_DB_HOST" default:"bills-psql"`
	DBPort        string `envconfig:"PIGLET_BILLS_DB_PORT" default:"5432"`
	DBName        string `envconfig:"PIGLET_BILLS_DB_NAME" default:"Accounting"`
}

type RabbitMQConfig struct {
	Server      string `envconfig:"RABBITMQ_SERVER" default:"piglet-rabbitmq"`
	Port        string `envconfig:"RABBITMQ_PORT" default:"5672"`
	RMQUser     string `envconfig:"RABBITMQ_USER" default:"user"`
	RMQPassword string `envconfig:"RABBITMQ_PASSWORD" default:"password"`
}

// InitConfig reads config variables from env and init *Config value
func MustLoadConfig() *Config {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		panic("piglet-bills: failed to read config: " + err.Error())
	}

	return cfg
}
