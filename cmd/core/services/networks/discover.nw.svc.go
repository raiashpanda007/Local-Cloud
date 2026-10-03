package service_networks

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

func DiscoverWorkers(serviceProtocol, domainName string, centralCtx, errorCtx context.Context) error {

	resolver, err := zeroconf.NewResolver(nil)

	if err != nil {
		return err
	}

	entriesRecieved := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Printf("Found Service: %s", entry.Instance)
			log.Printf("  Hostname: %s", entry.HostName)
			log.Printf("  IP Targets: %v", entry.AddrIPv4)
			log.Printf("  Port: %d", entry.Port)
			log.Printf("  Metadata (TXT): %v\n", entry.Text)
		}
	}(entriesRecieved)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Printf("Browsing for %s services...", serviceProtocol)
	err = resolver.Browse(ctx, serviceProtocol, domainName, entriesRecieved)
	if err != nil {
		log.Fatalf("Failed to browse: %v", err)
	}

	// Wait for the context timeout to expire
	<-ctx.Done()
	log.Println("Discovery browser closed.")

	return nil
}
