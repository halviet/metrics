package main

import (
	"github.com/halviet/metrics/internal/agent"
	"log"
	"time"
)

func main() {
	a := agent.New()

	if err := a.Start(2*time.Second, 10*time.Second); err != nil {
		log.Fatal(err)
	}
}
