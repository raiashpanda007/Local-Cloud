package core

import (
	"context"
	"os"

	"raiashpanda007/local-cloud-cli/cmd/core/env"
	service_networks "raiashpanda007/local-cloud-cli/cmd/core/services/networks"
	"raiashpanda007/local-cloud-cli/cmd/core/types"
)

const (
	DOMAIN_NAME string = "local.cloud"
)

// Rightnow calling this temp daemon.

func Daemon(centralCtx context.Context, nodeType types.NODE_TYPE) {
	var serviceName string
	var serviceProtocol string = "_http._tcp"
	metaData := []string{"version=1.0", "backend=golang"}

	env.Log.Info("daemon starting", "node", nodeType, "env", env.Env)
	// centralErrChan := make(chan error)

	errorCtx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	// Publish to mdns
	if nodeType == types.MASTER_NODE_TYPE {
		serviceName = "master"
	}

	env.Log.Debug("publishing service", "service", serviceName, "protocol", serviceProtocol, "domain", DOMAIN_NAME)

	err := service_networks.PublishServicemDNS(serviceName, serviceProtocol, DOMAIN_NAME, metaData, centralCtx, errorCtx)
	if err != nil {
		env.Log.Error("daemon failed", "err", err)
		os.Exit(1)
	}

	env.Log.Info("daemon stopped", "node", nodeType)
}
