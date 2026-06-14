package discovery

import (
	"context"
	"fmt"
	"log"

	"github.com/grandcat/zeroconf"
)

type Discovery struct {
	Name        string
	Port        int
	ServiceType string
	Domain      string

	OnPeerFound func(name, host string, port int)

	server   *zeroconf.Server
	resolver *zeroconf.Resolver
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(name string, port int) *Discovery {
	ctx, cancel := context.WithCancel(context.Background())
	return &Discovery{
		Name:        name,
		Port:        port,
		ServiceType: "_p2psync._tcp",
		Domain:      "local",
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (d *Discovery) Start() error {
	server, err := zeroconf.Register(d.Name, d.ServiceType, d.Domain, d.Port, nil, nil)
	if err != nil {
		return fmt.Errorf("zeroconf register: %w", err)
	}
	d.server = server

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("zeroconf resolver: %w", err)
	}
	d.resolver = resolver

	entries := make(chan *zeroconf.ServiceEntry)
	go d.handleEntries(entries)

	go func() {
		if err := resolver.Browse(d.ctx, d.ServiceType, d.Domain, entries); err != nil {
			log.Printf("zeroconf browse error: %v", err)
		}
	}()

	return nil
}

func (d *Discovery) handleEntries(entries chan *zeroconf.ServiceEntry) {
	for entry := range entries {
		if entry.Instance == d.Name {
			continue // skip self
		}
		host := ""
		if len(entry.AddrIPv4) > 0 {
			host = entry.AddrIPv4[0].String()
		} else if len(entry.AddrIPv6) > 0 {
			host = entry.AddrIPv6[0].String()
		}
		if host == "" && entry.HostName != "" {
			host = entry.HostName
		}
		if d.OnPeerFound != nil && host != "" {
			d.OnPeerFound(entry.Instance, host, entry.Port)
		}
	}
}

func (d *Discovery) Stop() {
	d.cancel()
	if d.server != nil {
		d.server.Shutdown()
	}
}
