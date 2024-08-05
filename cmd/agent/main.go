package main

import (
	"flag"
	"github.com/halviet/metrics/internal/agent"
	"log"
	"time"
)

func main() {
	srvAddr := flag.String("a", "localhost:8080", "address of a metrics server (addr:port)")
	reportInterval := flag.Int64("r", 10, "report interval, in seconds")
	pollInterval := flag.Int64("p", 2, "poll interval, in seconds")

	flag.Parse()

	a := agent.New()
	a.SrvAddr = *srvAddr

	if err := a.Start(time.Duration(*pollInterval)*time.Second, time.Duration(*reportInterval)*time.Second); err != nil {
		log.Fatal(err)
	}
}
