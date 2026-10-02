package service_networks

import (
	"context"
	"net"

	"raiashpanda007/local-cloud-cli/cmd/core/env"
	service_nw_server "raiashpanda007/local-cloud-cli/cmd/core/services/networks/server"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)

func PublishServicemDNS(serviceName, serviceProtocol, domainName string, metadata []string, centralCtx, errorCtx context.Context) error {
	listener, err := getRandomPort()
	if err != nil {
		env.Log.Error("tcp listen failed", "err", err)
		return err
	}

	serviceName += uuid.New().String()[:8]
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	env.Log.Debug("tcp listener open", "addr", listener.Addr().String())

	mdns, err := zeroconf.Register(serviceName, serviceProtocol, domainName, port, metadata, nil)
	if err != nil {
		env.Log.Error("mdns register failed", "service", serviceName, "port", port, "err", err)
		return err
	}
	defer mdns.Shutdown()

	env.Log.Info("mdns service published", "service", serviceName, "protocol", serviceProtocol, "domain", domainName, "port", port)

	return service_nw_server.WorkerServer(listener, centralCtx, errorCtx)

}
