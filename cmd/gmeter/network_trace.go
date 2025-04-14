package main

import "github.com/spf13/cobra"

var traceCmd = &cobra.Command{
	Use:   "network-trace",
	Short: "trace the network packets flowing into a pid",
}

var pid uint32

func init() {
	counterCmd.Flags().Uint32VarP(&pid, "pid", "p", 0, "PID to monitor")
	rootCmd.AddCommand(counterCmd)
}

func runNetworkTrace(cmd *cobra.Command, args []string) {

}
