package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "localcloud",
	Short: "Pool computing resources across your machines",
	Long: `Local Cloud connects computers you control and combines their computing
resources into a single pool.

A master coordinates connected machines, manages workers, and distributes
workloads. Workers run on participating machines and execute assigned tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
