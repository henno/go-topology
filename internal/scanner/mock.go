package scanner

import (
	"time"
)

// MockScanner returns fixture devices for testing
type MockScanner struct{}

// NewMockScanner creates a new mock scanner instance
func NewMockScanner() *MockScanner {
	return &MockScanner{}
}

// Scan returns 5 predefined fixture devices
func (m *MockScanner) Scan(network string) (<-chan Device, error) {
	ch := make(chan Device, 5)

	// Return fixture devices in a goroutine to simulate async discovery
	go func() {
		defer close(ch)

		fixtures := []Device{
			{
				IP:       "192.168.1.1",
				MAC:      "00:0C:42:12:34:56",
				Hostname: "gateway.local",
				Vendor:   "Routerboard.com",
				Type:     "router",
			},
			{
				IP:       "192.168.1.2",
				MAC:      "00:17:88:AB:CD:EF",
				Hostname: "switch-01",
				Vendor:   "Zyxel",
				Type:     "switch",
			},
			{
				IP:       "192.168.1.10",
				MAC:      "D8:9E:F3:11:22:33",
				Hostname: "workstation-01",
				Vendor:   "Dell",
				Type:     "computer",
			},
			{
				IP:       "192.168.1.20",
				MAC:      "00:1E:C9:44:55:66",
				Hostname: "printer-01",
				Vendor:   "HP",
				Type:     "printer",
			},
			{
				IP:       "192.168.1.99",
				MAC:      "AA:BB:CC:DD:EE:FF",
				Hostname: "",
				Vendor:   "Unknown",
				Type:     "unknown",
			},
		}

		// Emit devices with a small delay to simulate discovery
		for _, device := range fixtures {
			ch <- device
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return ch, nil
}

