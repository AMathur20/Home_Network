package models

import "time"

type DeviceType string

const (
	DeviceTypeMikroTik   DeviceType = "mikrotik"
	DeviceTypeUniFi      DeviceType = "unifi"
	DeviceTypeEdgeRouter DeviceType = "edgerouter"
	DeviceTypeGeneric    DeviceType = "generic"
)

type Config struct {
	Poller IntervalConfig `yaml:"poller"`
	Devices []DeviceConfig `yaml:"devices"`
}

type IntervalConfig struct {
	Live     int `yaml:"live"`     // seconds
	History  int `yaml:"history"`  // seconds
}

type DeviceConfig struct {
	Name     string     `yaml:"name"`
	Host     string     `yaml:"host"`
	Type     DeviceType `yaml:"type"`
	Auth     AuthConfig `yaml:"auth"`
	SNMP     SNMPConfig `yaml:"snmp"`
}

type AuthConfig struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type SNMPConfig struct {
	Version   string `yaml:"version"` // v2c, v3
	Community string `yaml:"community,omitempty"`
	Port      int    `yaml:"port"`
	// v3 fields can be added here
}

type InterfaceMetric struct {
	DeviceName    string
	InterfaceName string
	Timestamp     time.Time
	InOctets      uint64
	OutOctets     uint64
	InSpeed       float64 // bps
	OutSpeed      float64 // bps
	Status        string  // up, down
}
