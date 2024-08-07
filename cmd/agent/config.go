package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	SrvAddr        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
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
	if envs.ReportInterval != 0 {
		cfg.ReportInterval = envs.ReportInterval
	}
	if envs.PollInterval != 0 {
		cfg.PollInterval = envs.PollInterval
	}

	return cfg, nil
}

func parseFlags() Config {
	srvAddr := flag.String("a", "localhost:8080", "address of a metrics server (addr:port)")
	reportInterval := flag.Int64("r", 10, "report interval, in seconds")
	pollInterval := flag.Int64("p", 2, "poll interval, in seconds")

	flag.Parse()

	return Config{
		SrvAddr:        *srvAddr,
		PollInterval:   time.Duration(*pollInterval) * time.Second,
		ReportInterval: time.Duration(*reportInterval) * time.Second,
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
