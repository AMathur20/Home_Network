package topology

import (
	"fmt"
	"log"
	"time"

	"github.com/AMathur20/Home_Network/internal/models"
	"github.com/gosnmp/gosnmp"
)

const (
	oidLldpRemSysName = ".1.0.8802.1.1.2.1.4.1.1.9"
	oidLldpRemPortId  = ".1.0.8802.1.1.2.1.4.1.1.7"
	oidIfName         = ".1.3.6.1.2.1.31.1.1.1.1"
)

type Crawler struct {
	devices []models.DeviceConfig
}

func NewCrawler(devices []models.DeviceConfig) *Crawler {
	return &Crawler{devices: devices}
}

func (c *Crawler) Discover() (*Topology, error) {
	log.Println("Starting real LLDP topology discovery...")

	links := make([]Link, 0)

	for _, dev := range c.devices {
		deviceLinks, err := c.discoverForDevice(dev)
		if err != nil {
			log.Printf("Error discovering neighbors for %s: %v", dev.Name, err)
			continue
		}
		links = append(links, deviceLinks...)
	}

	return &Topology{Links: links}, nil
}

func (c *Crawler) discoverForDevice(dev models.DeviceConfig) ([]Link, error) {
	params := &gosnmp.GoSNMP{
		Target:    dev.Host,
		Port:      uint16(dev.SNMP.Port),
		Community: dev.SNMP.Community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   3,
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}
	defer params.Conn.Close()

	// 1. Map ifIndex to ifName
	ifNames := make(map[int]string)
	err = params.BulkWalk(oidIfName, func(pdu gosnmp.SnmpPDU) error {
		index := 0
		fmt.Sscanf(pdu.Name, oidIfName+".%d", &index)
		ifNames[index] = models.PduToString(pdu.Value)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk ifName: %v", err)
	}

	links := make([]Link, 0)

	// 2. Discover neighbors
	// The LLDP Rem Table index is roughly: lldpRemTimeMark.lldpLocalPortNum.lldpRemIndex
	err = params.BulkWalk(oidLldpRemSysName, func(pdu gosnmp.SnmpPDU) error {
		// Extract local port index from OID
		// OID format: .1.0.8802.1.1.2.1.4.1.1.9.<timeMark>.<localPortNum>.<remIndex>
		suffix := pdu.Name[len(oidLldpRemSysName)+1:]
		var timeMark, localPortNum, remIndex int
		fmt.Sscanf(suffix, "%d.%d.%d", &timeMark, &localPortNum, &remIndex)

		sourceIface := ifNames[localPortNum]
		if sourceIface == "" {
			sourceIface = fmt.Sprintf("port-%d", localPortNum)
		}

		targetDevice := models.PduToString(pdu.Value)

		// Fetch the remote port ID for this specific entry
		remotePortPath := fmt.Sprintf("%s.%d.%d.%d", oidLldpRemPortId, timeMark, localPortNum, remIndex)
		result, err := params.Get([]string{remotePortPath})
		targetIface := "unknown"
		if err == nil && len(result.Variables) > 0 {
			targetIface = models.PduToString(result.Variables[0].Value)
		}

		links = append(links, Link{
			SourceDevice:    dev.Name,
			SourceInterface: sourceIface,
			TargetDevice:    targetDevice,
			TargetInterface: targetIface,
			Type:            ClassifyLink(sourceIface),
		})

		return nil
	})

	return links, err
}
