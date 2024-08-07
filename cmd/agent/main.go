package main

import (
	"github.com/halviet/metrics/internal/agent"
	"log"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("config initialization fail: %v", err)
	}

	a := agent.New()
	a.SrvAddr = cfg.SrvAddr

	if err := a.Start(cfg.PollInterval, cfg.ReportInterval); err != nil {
		log.Fatal(err)
	}
}
