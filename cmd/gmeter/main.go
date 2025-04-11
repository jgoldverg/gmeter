package main

import (
	"flag"
	"gmeter/pkg/counter"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	iface := flag.String("i", "eno1np0:", "Network interface")
	flag.Parse()

	counter, err := counter.NewCounter()
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}
	defer counter.Close()

	if err := counter.Attach(*iface); err != nil {
		log.Fatalf("Failed to attach: %v", err)
	}

	log.Printf("Counting packets on %s...", *iface)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-tick.C:
			count, err := counter.ReadCount()
			if err != nil {
				log.Printf("Error reading count: %v", err)
				continue
			}
			log.Printf("Packets: %d", count)
		case <-sig:
			return
		}
	}
}
