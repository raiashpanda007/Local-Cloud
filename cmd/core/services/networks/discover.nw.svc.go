package service_networks

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"raiashpanda007/local-cloud-cli/cmd/core/env"

	"github.com/grandcat/zeroconf"
)

func DiscoverWorkers(serviceProtocol, domainName string, centralCtx, errorCtx context.Context) error {
	resolver, err := zeroconf.NewResolver()
	if err != nil {
		env.Log.Error("mdns resolver failed", "err", err)
		return err
	}

	ctx, cancel := context.WithCancel(centralCtx)
	defer cancel()

	go func() {
		select {
		case <-errorCtx.Done():
			cancel()
		case <-ctx.Done():
		}
	}()

	entries := make(chan *zeroconf.ServiceEntry)
	if err := resolver.Browse(ctx, serviceProtocol, domainName, entries); err != nil {
		env.Log.Error("mdns browse failed", "protocol", serviceProtocol, "domain", domainName, "err", err)
		return err
	}

	env.Log.Info("browsing for workers", "protocol", serviceProtocol, "domain", domainName)

	caught := map[string]string{}
	for entry := range entries {
		if !strings.HasPrefix(entry.Instance, "worker") {
			env.Log.Debug("ignored service", "instance", entry.Instance)
			continue
		}
		if addr, ok := caught[entry.Instance]; ok {
			env.Log.Debug("worker already caught", "instance", entry.Instance, "addr", addr)
			continue
		}

		env.Log.Debug("worker found", "instance", entry.Instance, "host", entry.HostName, "port", entry.Port)

		addr, err := catchWorker(ctx, entry)
		if err != nil {
			env.Log.Error("worker catch failed", "instance", entry.Instance, "err", err)
			continue
		}

		caught[entry.Instance] = addr
		env.Log.Info("worker caught", "instance", entry.Instance, "addr", addr, "workers", len(caught))
	}

	env.Log.Info("stopped browsing for workers", "workers", len(caught))
	return nil
}

func catchWorker(ctx context.Context, entry *zeroconf.ServiceEntry) (string, error) {
	ips := append([]net.IP{}, entry.AddrIPv4...)
	ips = append(ips, entry.AddrIPv6...)
	if len(ips) == 0 {
		return "", fmt.Errorf("worker %s has no address", entry.Instance)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	var lastErr error

	for _, ip := range ips {
		addr := net.JoinHostPort(ip.String(), strconv.Itoa(entry.Port))
		reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "http://"+addr+"/healthy", nil)
		if err != nil {
			cancel()
			lastErr = err
			continue
		}

		resp, err := client.Do(req)
		cancel()
		if err != nil {
			env.Log.Debug("worker health check failed", "instance", entry.Instance, "addr", addr, "err", err)
			lastErr = err
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return addr, nil
		}
		lastErr = fmt.Errorf("worker %s at %s returned %s", entry.Instance, addr, resp.Status)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("worker %s could not be reached", entry.Instance)
	}
	return "", lastErr
}
