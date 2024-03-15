package config

import (
	"flag"
	"os"
)

type Config struct {
	Port string
	Env  string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() error {
	flag.StringVar(&cfg.Port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENV"), "Environment (development|staging|production)")

	return nil
}
