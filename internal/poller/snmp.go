package poller

import (
	"fmt"
	"log"
	"time"

	"github.com/AMathur20/Home_Network/internal/models"
	"github.com/AMathur20/Home_Network/internal/snmphelper"
	"github.com/gosnmp/gosnmp"
)

const (
	oidIfName        = ".1.3.6.1.2.1.31.1.1.1.1"
	oidIfHCInOctets  = ".1.3.6.1.2.1.31.1.1.1.6"
	oidIfHCOutOctets = ".1.3.6.1.2.1.31.1.1.1.10"
	oidIfInOctets    = ".1.3.6.1.2.1.2.2.1.10"
	oidIfOutOctets   = ".1.3.6.1.2.1.2.2.1.16"
	oidIfOperStatus  = ".1.3.6.1.2.1.2.2.1.8"
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

	timestamp := time.Now()
	metrics := make(map[int]*models.InterfaceMetric)

	// 1. Fetch Interface Names
	err = params.BulkWalk(oidIfName, func(pdu gosnmp.SnmpPDU) error {
		index := 0
		fmt.Sscanf(pdu.Name, oidIfName+".%d", &index)
		metrics[index] = &models.InterfaceMetric{
			DeviceName:    p.config.Name,
			InterfaceName: snmphelper.PduToString(pdu.Value),
			Timestamp:     timestamp,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk ifName: %v", err)
	}

	// 2. Fetch Operational Status
	err = params.BulkWalk(oidIfOperStatus, func(pdu gosnmp.SnmpPDU) error {
		index := 0
		fmt.Sscanf(pdu.Name, oidIfOperStatus+".%d", &index)
		if m, ok := metrics[index]; ok {
			status := snmphelper.PduToInt(pdu.Value)
			if status == 1 {
				m.Status = "up"
			} else {
				m.Status = "down"
			}
		}
		return nil
	})

	// 3. Fetch Counters (Prefer 64-bit HC counters)
	useHC := true
	err = params.BulkWalk(oidIfHCInOctets, func(pdu gosnmp.SnmpPDU) error {
		index := 0
		fmt.Sscanf(pdu.Name, oidIfHCInOctets+".%d", &index)
		if m, ok := metrics[index]; ok {
			m.InOctets = snmphelper.PduToUint64(pdu.Value)
		}
		return nil
	})
	if err != nil {
		log.Printf("Device %s does not support ifHCInOctets, falling back to 32-bit", p.config.Name)
		useHC = false
	}

	if useHC {
		params.BulkWalk(oidIfHCOutOctets, func(pdu gosnmp.SnmpPDU) error {
			index := 0
			fmt.Sscanf(pdu.Name, oidIfHCOutOctets+".%d", &index)
			if m, ok := metrics[index]; ok {
				m.OutOctets = snmphelper.PduToUint64(pdu.Value)
			}
			return nil
		})
	} else {
		// Fallback to 32-bit counters
		params.BulkWalk(oidIfInOctets, func(pdu gosnmp.SnmpPDU) error {
			index := 0
			fmt.Sscanf(pdu.Name, oidIfInOctets+".%d", &index)
			if m, ok := metrics[index]; ok {
				m.InOctets = uint64(snmphelper.PduToUint64(pdu.Value))
			}
			return nil
		})
		params.BulkWalk(oidIfOutOctets, func(pdu gosnmp.SnmpPDU) error {
			index := 0
			fmt.Sscanf(pdu.Name, oidIfOutOctets+".%d", &index)
			if m, ok := metrics[index]; ok {
				m.OutOctets = uint64(snmphelper.PduToUint64(pdu.Value))
			}
			return nil
		})
	}

	result := make([]models.InterfaceMetric, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, *m)
	}

	return result, nil
}
