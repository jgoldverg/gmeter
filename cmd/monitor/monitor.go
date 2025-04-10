package monitor

import (
	"fmt"
	"github.com/cilium/ebpf"
	"github.com/spf13/cobra"
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "A Cli app using Cillium eBPF",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startMonitoring(); err != nil {
			fmt.Println("Monitoring failed: %v\n", err)
			return
		}
		fmt.Println("Monitoring started...")
	},
}

func startMonitoring() error {
	coll, err := ebpf.LoadCollection("build/bpf/monitor.o")
	if err != nil {
		return fmt.Errorf("loading collection: %s", err)
	}

	defer coll.Close()

	return nil
}
