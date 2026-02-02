package poller

import (
	"log"
	"time"

	"github.com/AMathur20/Home_Network/internal/models"
	"github.com/AMathur20/Home_Network/internal/storage"
	"github.com/AMathur20/Home_Network/internal/topology"
)

type PollingEngine struct {
	config   *models.Config
	storage  *storage.DuckDBStorage
	topology *topology.Topology
	state    map[string]interfaceState
}

type interfaceState struct {
	lastInOctets  uint64
	lastOutOctets uint64
	lastTime      time.Time
}

func NewPollingEngine(cfg *models.Config, s *storage.DuckDBStorage, t *topology.Topology) *PollingEngine {
	return &PollingEngine{
		config:   cfg,
		storage:  s,
		topology: t,
		state:    make(map[string]interfaceState),
	}
}

func (e *PollingEngine) ReloadTopology(t *topology.Topology) {
	e.topology = t
	log.Println("Engine topology reloaded.")
}

func (e *PollingEngine) Start() {
	log.Printf("Starting polling engine with %d devices", len(e.config.Devices))

	ticker := time.NewTicker(time.Duration(e.config.Poller.Live) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, dev := range e.config.Devices {
			go e.pollDevice(dev)
		}
	}
}

func (e *PollingEngine) pollDevice(dev models.DeviceConfig) {
	var metrics []models.InterfaceMetric
	var err error

	if dev.Type == models.DeviceTypeMikroTik || dev.SNMP.Version != "" {
		p := NewSNMPPoller(dev)
		metrics, err = p.Poll()
	}

	if err != nil {
		log.Printf("Error polling device %s: %v", dev.Name, err)
		return
	}

	for _, m := range metrics {
		key := m.DeviceName + "/" + m.InterfaceName
		last, ok := e.state[key]
		if ok {
			duration := m.Timestamp.Sub(last.lastTime).Seconds()
			if duration > 0 {
				// Bps = (delta octets * 8) / seconds
				if m.InOctets >= last.lastInOctets {
					m.InSpeed = float64(m.InOctets-last.lastInOctets) * 8 / duration
				}
				if m.OutOctets >= last.lastOutOctets {
					m.OutSpeed = float64(m.OutOctets-last.lastOutOctets) * 8 / duration
				}
			}
		}
		e.state[key] = interfaceState{
			lastInOctets:  m.InOctets,
			lastOutOctets: m.OutOctets,
			lastTime:      m.Timestamp,
		}

		if err := e.storage.SaveMetric(m); err != nil {
			log.Printf("Error saving metric for %s/%s: %v", m.DeviceName, m.InterfaceName, err)
		}
	}
}
