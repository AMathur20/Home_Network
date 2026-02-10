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

	// MikroTik MNDP OIDs
	// mtxrNeighborTable: .1.3.6.1.4.1.14988.1.1.11.1
	oidMndpNeighborIdentity  = ".1.3.6.1.4.1.14988.1.1.11.1.3" // Neighbor System Name
	oidMndpNeighborInterface = ".1.3.6.1.4.1.14988.1.1.11.1.8" // Local Interface on which neighbor was discovered
)

type Crawler struct {
	devices []models.DeviceConfig
}

func NewCrawler(devices []models.DeviceConfig) *Crawler {
	return &Crawler{devices: devices}
}

func (c *Crawler) Discover() (*Topology, error) {
	log.Println("Starting topology discovery (LLDP + MNDP)...")

	links := make([]Link, 0)

	for _, dev := range c.devices {
		// LLDP discovery (Universal)
		lldpLinks, err := c.discoverLldpForDevice(dev)
		if err != nil {
			log.Printf("Error discovering LLDP neighbors for %s: %v", dev.Name, err)
		} else {
			links = append(links, lldpLinks...)
		}

		// MNDP discovery (MikroTik specific)
		if dev.Type == models.DeviceTypeMikroTik {
			mndpLinks, err := c.discoverMndpForDevice(dev)
			if err != nil {
				log.Printf("Error discovering MNDP neighbors for %s: %v", dev.Name, err)
			} else {
				links = append(links, mndpLinks...)
			}
		}
	}

	return &Topology{Links: deduplicateLinks(links)}, nil
}

func (c *Crawler) discoverLldpForDevice(dev models.DeviceConfig) ([]Link, error) {
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

func (c *Crawler) discoverMndpForDevice(dev models.DeviceConfig) ([]Link, error) {
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

	links := make([]Link, 0)

	err = params.BulkWalk(oidMndpNeighborIdentity, func(pdu gosnmp.SnmpPDU) error {
		// Extract index from OID suffix
		suffix := pdu.Name[len(oidMndpNeighborIdentity)+1:]
		targetDevice := models.PduToString(pdu.Value)

		// Fetch the local interface for this specific neighbor entry
		localIfacePath := fmt.Sprintf("%s.%s", oidMndpNeighborInterface, suffix)
		result, err := params.Get([]string{localIfacePath})
		sourceIface := "unknown"
		if err == nil && len(result.Variables) > 0 {
			sourceIface = models.PduToString(result.Variables[0].Value)
		}

		links = append(links, Link{
			SourceDevice:    dev.Name,
			SourceInterface: sourceIface,
			TargetDevice:    targetDevice,
			TargetInterface: "unknown", // MNDP often doesn't provide remote port ID via SNMP easily
			Type:            ClassifyLink(sourceIface),
		})

		return nil
	})

	return links, err
}

func deduplicateLinks(links []Link) []Link {
	type key struct {
		srcDev, srcIf, tgtDev string
	}
	seen := make(map[key]bool)
	unique := make([]Link, 0)

	for _, l := range links {
		k := key{l.SourceDevice, l.SourceInterface, l.TargetDevice}
		if !seen[k] {
			seen[k] = true
			unique = append(unique, l)
		}
	}
	return unique
}
