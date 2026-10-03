package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the master on this machine",
	Long: `Start runs the master on this machine.

The master coordinates connected machines, discovers workers, manages
connections, and distributes workloads.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("master does not run inside the cli")
	},
}

func init() {
	masterCmd.AddCommand(startCmd)
}
