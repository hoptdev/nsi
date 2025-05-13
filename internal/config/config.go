package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string       `yaml:"env"`
	PSQL_Connect string       `yaml:"psql_connect"`
	Server       ServerConfig `yaml:"server"`
	Client       ClientConfig `yaml:"client"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Address string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}

type ClientConfig struct {
	Port    int    `yaml:"port"`
	Address string `yaml:"addr"`
}

func Load() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("read config failed")
	}

	return &cfg
}

func fetchConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}
