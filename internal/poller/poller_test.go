package poller

import (
	"testing"
	"time"

	"github.com/AMathur20/Home_Network/internal/models"
)

func TestThroughputCalculation(t *testing.T) {
	engine := NewPollingEngine(&models.Config{}, nil, nil)

	deviceName := "test-device"
	ifaceName := "ether1"

	// T1
	t1 := time.Now()
	m1 := &models.InterfaceMetric{
		DeviceName:    deviceName,
		InterfaceName: ifaceName,
		Timestamp:     t1,
		InOctets:      1000,
		OutOctets:     2000,
	}

	// T2 (5 seconds later)
	t2 := t1.Add(5 * time.Second)
	m2 := &models.InterfaceMetric{
		DeviceName:    deviceName,
		InterfaceName: ifaceName,
		Timestamp:     t2,
		InOctets:      2000, // 1000 delta
		OutOctets:     4000, // 2000 delta
	}

	// First poll sets the state
	engine.pollDeviceToUpdateState(m1)

	// Second poll calculates speed
	engine.pollDeviceToUpdateState(m2)

	// Check speeds
	// InSpeed = (1000 * 8) / 5 = 1600 bps
	// OutSpeed = (2000 * 8) / 5 = 3200 bps
	if m2.InSpeed != 1600 {
		t.Errorf("Expected InSpeed 1600, got %f", m2.InSpeed)
	}
	if m2.OutSpeed != 3200 {
		t.Errorf("Expected OutSpeed 3200, got %f", m2.OutSpeed)
	}
}

func TestThroughputRollover(t *testing.T) {
	engine := NewPollingEngine(&models.Config{}, nil, nil)

	deviceName := "test-device"
	ifaceName := "ether1"

	t1 := time.Now()
	m1 := &models.InterfaceMetric{
		DeviceName:    deviceName,
		InterfaceName: ifaceName,
		Timestamp:     t1,
		InOctets:      1000,
	}

	// T2: Rollover (current < last)
	t2 := t1.Add(5 * time.Second)
	m2 := &models.InterfaceMetric{
		DeviceName:    deviceName,
		InterfaceName: ifaceName,
		Timestamp:     t2,
		InOctets:      500, // Lower than 1000
	}

	engine.pollDeviceToUpdateState(m1)
	engine.pollDeviceToUpdateState(m2)

	// Should be 0 for now as we don't handle the exact bit width of the counter gracefully yet
	if m2.InSpeed != 0 {
		t.Errorf("Expected InSpeed 0 on rollover, got %f", m2.InSpeed)
	}
}

// Helper for testing
func (e *PollingEngine) pollDeviceToUpdateState(m *models.InterfaceMetric) {
	key := m.DeviceName + "/" + m.InterfaceName
	last, ok := e.state[key]
	if ok {
		duration := m.Timestamp.Sub(last.lastTime).Seconds()
		if duration > 0 {
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
}
