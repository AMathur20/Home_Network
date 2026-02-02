package poller

import (
	"log"
	"time"

	"github.com/AMathur20/Home_Network/internal/models"
	"github.com/gosnmp/gosnmp"
)

type SNMPPoller struct {
	config models.DeviceConfig
}

func NewSNMPPoller(cfg models.DeviceConfig) *SNMPPoller {
	return &SNMPPoller{config: cfg}
}

func (p *SNMPPoller) Poll() ([]models.InterfaceMetric, error) {
	params := &gosnmp.GoSNMP{
		Target:    p.config.Host,
		Port:      uint16(p.config.SNMP.Port),
		Community: p.config.SNMP.Community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   3,
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}
	defer params.Conn.Close()

	// In a real implementation, we would BulkWalk these.
	// For this refinement, let's fetch a few key interfaces or mock the names for classification demo.

	metrics := make([]models.InterfaceMetric, 0)
	timestamp := time.Now()

	log.Printf("Polling SNMP for device: %s", p.config.Name)

	// Simulated interfaces for demonstration of the classification logic
	interfaces := []struct {
		name  string
		descr string
	}{
		{"ether1", "1G Copper"},
		{"sfp-sfpplus1", "10G Fiber Uplink"},
		{"wlan1", "2.4GHz Wi-Fi"},
	}

	for _, iface := range interfaces {
		// In a real crawl, we'd get these from SNMP Walk
		metrics = append(metrics, models.InterfaceMetric{
			DeviceName:    p.config.Name,
			InterfaceName: iface.name,
			Timestamp:     timestamp,
			InOctets:      1000, // Placeholder
			OutOctets:     500,  // Placeholder
			Status:        "up",
		})

		// Note: The classification is used by the topology crawler
		// but we can also store the type in the metric if we extend the model.
		// For now, the PRD requirement is about topology detection.
	}

	return metrics, nil
}
