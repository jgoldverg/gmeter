package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"gmeter/gen"
)

func main() {
	// Parse CLI args
	ifaceName := flag.String("i", "eth0", "Network interface to monitor")
	flag.Parse()

	// Remove memory limits
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove memlock limit: %v", err)
	}

	// Load eBPF objects

	var objs gen.CounterObject
	if err := gen.loadCounter(&objs, nil); err != nil {
		log.Fatalf("Loading BPF objects: %v", err)
	}
	defer objs.Close()

	// Get network interface
	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatalf("Failed to get interface %s: %v", *ifaceName, err)
	}

	// Attach XDP program
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.CountPackets,
		Interface: iface.Index,
		Flags:     link.XDPGenericMode,
	})
	if err != nil {
		log.Fatalf("Failed to attach XDP program: %v", err)
	}
	defer l.Close()

	log.Printf("Counting packets on %s (Ctrl+C to exit)", *ifaceName)

	// Setup monitoring
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-tick.C:
			var count uint64
			if err := objs.PktCount.Lookup(uint32(0), &count); err != nil {
				log.Printf("Error reading counter: %v", err)
				continue
			}
			log.Printf("Packets received: %d", count)
		case <-stop:
			log.Println("Detaching and exiting...")
			return
		}
	}
}
