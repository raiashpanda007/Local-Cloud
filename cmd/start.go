package cmd

import (
	"github.com/spf13/cobra"
	"raiashpanda007/local-cloud-cli/cmd/core"
	"raiashpanda007/local-cloud-cli/cmd/core/types"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the master on this machine",
	Long: `Start runs the master on this machine.

The master coordinates connected machines, discovers workers, manages
connections, and distributes workloads.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		core.Daemon(cmd.Context(), types.MASTER_NODE_TYPE)
		return nil
	},
}

func init() {
	masterCmd.AddCommand(startCmd)
}
