package topology

import (
	"strings"
)

type LinkType string

const (
	LinkType10G      LinkType = "10g"
	LinkType1G       LinkType = "1g"
	LinkTypeEthernet LinkType = "ethernet"
	LinkTypeWireless LinkType = "wireless"
)

type Link struct {
	SourceDevice    string   `yaml:"source_device"`
	SourceInterface string   `yaml:"source_interface"`
	TargetDevice    string   `yaml:"target_device"`
	TargetInterface string   `yaml:"target_interface"`
	Type            LinkType `yaml:"type"`
	Manual          bool     `yaml:"manual,omitempty"` // True if manually added/overridden
}

type Topology struct {
	Links []Link `yaml:"links"`
}

func ClassifyLink(ifDescr string) LinkType {
	lower := strings.ToLower(ifDescr)
	switch {
	case strings.Contains(lower, "sfpplus"), strings.Contains(lower, "10g"):
		return LinkType10G
	case strings.Contains(lower, "sfp"):
		return LinkType1G
	case strings.Contains(lower, "wlan"), strings.Contains(lower, "wifi"):
		return LinkTypeWireless
	default:
		return LinkTypeEthernet
	}
}
