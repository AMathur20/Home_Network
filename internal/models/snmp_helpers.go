package models

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
)

// PduToString safely converts an SNMP PDU value to a string.
func PduToString(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// PduToInt safely converts an SNMP PDU value to an int.
func PduToInt(val interface{}) int {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case uint:
		return int(v)
	case uint64:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}

// PduToUint64 safely converts an SNMP PDU value to a uint64.
func PduToUint64(val interface{}) uint64 {
	if val == nil {
		return 0
	}
	return gosnmp.ToBigInt(val).Uint64()
}
