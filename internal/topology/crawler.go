package topology

import (
	"log"

	"github.com/AMathur20/Home_Network/internal/models"
)

type Crawler struct {
	devices []models.DeviceConfig
}

func NewCrawler(devices []models.DeviceConfig) *Crawler {
	return &Crawler{devices: devices}
}

func (c *Crawler) Discover() (*Topology, error) {
	log.Println("Starting LLDP topology discovery...")

	links := make([]Link, 0)

	// In a full implementation, we would:
	// 1. Walk lldpRemSysName (.1.0.8802.1.1.2.1.4.1.1.9)
	// 2. Walk lldpRemPortId (.1.0.8802.1.1.2.1.4.1.1.7)
	// 3. Map local interface index to name using ifDescr

	// Simulated discovery for demonstration
	for _, dev := range c.devices {
		if dev.Name == "mikrotik-core" {
			links = append(links, Link{
				SourceDevice:    dev.Name,
				SourceInterface: "sfp-sfpplus1",
				TargetDevice:    "unifi-switch-agg",
				TargetInterface: "sfp-sfpplus1",
				Type:            ClassifyLink("sfp-sfpplus1"), // Auto-classified as 10G
			})
			links = append(links, Link{
				SourceDevice:    dev.Name,
				SourceInterface: "ether2",
				TargetDevice:    "living-room-sw",
				TargetInterface: "port1",
				Type:            ClassifyLink("ether2"), // Auto-classified as 1G/Ethernet
			})
		}
	}

	return &Topology{Links: links}, nil
}
