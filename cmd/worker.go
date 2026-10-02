package cmd

import (
	"raiashpanda007/local-cloud-cli/cmd/core"
	"raiashpanda007/local-cloud-cli/cmd/core/types"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Execute assigned tasks on this machine",
	Long: `The worker runs on a participating machine. It advertises availability
on the local network and executes tasks assigned by the master.

Remote connectivity is planned for a later release.`,
}

var workerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the worker on this machine",
	Long: `Start runs the worker on this machine.

The worker advertises availability on the local network and executes
tasks assigned by the master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		core.Daemon(cmd.Context(), types.WORKER_NODE_TYPE)
		return nil
	},
}

func init() {
	workerCmd.AddCommand(workerStartCmd)
	rootCmd.AddCommand(workerCmd)
}
