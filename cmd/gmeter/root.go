package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gmeter",
	Short: "eBPF powered networking observability tool",
}
