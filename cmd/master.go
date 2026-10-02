package cmd

import "github.com/spf13/cobra"

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Coordinate machines and distribute workloads",
	Long: `The master coordinates machines joined to a Local Cloud pool. It discovers
available workers, manages their connections, and distributes workloads
across them.

Workers advertise availability on the local network. Remote connectivity
is planned for a later release.`,
}

func init() {
	rootCmd.AddCommand(masterCmd)
}
