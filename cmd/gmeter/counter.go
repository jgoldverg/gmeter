package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gmeter/pkg/counter"
)

var iface string

var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Count packets on an interface",
	Run:   runCounter, // Reference to standalone function
}

func init() {
	counterCmd.Flags().StringVarP(&iface, "interface", "i", "eno1np0", "Network interface")
	rootCmd.AddCommand(counterCmd)
}

// runCounter contains all the counter logic
func runCounter(cmd *cobra.Command, args []string) {
	c, err := counter.NewCounter()
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}

	defer c.Close()

	if err := c.Attach(iface); err != nil {
		log.Fatalf("Failed to attach: %v", err)
	}

	log.Printf("Counting packets on %s...", iface)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-tick.C:
			count, err := c.ReadCount()
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
