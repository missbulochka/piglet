package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC        GRPCConfig    `json:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("piglet-auth: config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("piglet-auth: config path does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("piglet-auth: failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	// --config="path" or ENV
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	fmt.Println("1: ", res)

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	fmt.Println("2: ", res)

	return res
}
