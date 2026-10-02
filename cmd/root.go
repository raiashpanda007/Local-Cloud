package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"raiashpanda007/local-cloud-cli/cmd/core/env"

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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var finished atomic.Bool
	go func() {
		<-ctx.Done()
		if !finished.Load() {
			env.Log.Debug("shutdown signal received")
			os.Exit(130)
		}
	}()

	err := rootCmd.ExecuteContext(ctx)
	finished.Store(true)
	if err != nil {
		os.Exit(1)
	}
}
