package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "localcloud",
	Short: "Pool computing resources across your machines",
	Long: `Local Cloud connects computers you control and combines their computing
resources into a single pool.

A master coordinates connected machines, manages workers, and distributes
workloads. Workers run on participating machines, advertise availability
on the local network, and execute assigned tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	// This central context is the system wide context carrier.
	centralCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	go func() {
		<-centralCtx.Done()
		os.Exit(130)
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
