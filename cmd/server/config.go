package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	SrvAddr string `env:"ADDRESS"`
}

func NewConfig() (Config, error) {
	cfg := parseFlags()
	envs, err := parseEnv()
	if err != nil {
		return Config{}, err
	}

	if envs.SrvAddr != "" {
		cfg.SrvAddr = envs.SrvAddr
	}

	return cfg, nil
}

func parseFlags() Config {
	srvAddr := flag.String("a", "localhost:8080", "address of a metrics server (addr:port)")

	flag.Parse()

	return Config{
		SrvAddr: *srvAddr,
	}
}

func parseEnv() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parseEnv: %v", err)
	}

	return cfg, nil
}
