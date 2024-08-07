package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	SrvAddr        string
	PollInterval   time.Duration
	ReportInterval time.Duration
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
	var cfg struct {
		SrvAddr        string `env:"ADDRESS"`
		PollInterval   int    `env:"POLL_INTERVAL"`
		ReportInterval int    `env:"REPORT_INTERVAL"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parseEnv: %v", err)
	}

	return Config{
		SrvAddr:        cfg.SrvAddr,
		PollInterval:   time.Duration(cfg.PollInterval) * time.Second,
		ReportInterval: time.Duration(cfg.ReportInterval) * time.Second,
	}, nil
}
