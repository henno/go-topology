package scanner

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// RealScanner performs actual network scanning
type RealScanner struct {
	workers   int
	timeoutMs int
}

// NewRealScanner creates a new real scanner instance
func NewRealScanner(workers, timeoutMs int) *RealScanner {
	if workers <= 0 {
		workers = 100 // default
	}
	if timeoutMs <= 0 {
		timeoutMs = 1000 // default 1 second
	}
	return &RealScanner{
		workers:   workers,
		timeoutMs: timeoutMs,
	}
}

// Scan discovers devices on the network using ICMP ping
func (r *RealScanner) Scan(network string) (<-chan Device, error) {
	// Parse CIDR
	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, fmt.Errorf("invalid network CIDR: %w", err)
	}

	ch := make(chan Device, 100)

	// Generate list of IPs to scan
	ips := generateIPs(ipNet)

	// Start scanning in background
	go func() {
		defer close(ch)

		// Create work queue
		ipChan := make(chan string, len(ips))
		for _, ip := range ips {
			ipChan <- ip
		}
		close(ipChan)

		// Start workers
		var wg sync.WaitGroup
		for i := 0; i < r.workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ip := range ipChan {
					if device, ok := r.scanIP(ip); ok {
						ch <- device
					}
				}
			}()
		}

		wg.Wait()
	}()

	return ch, nil
}

// scanIP probes a single IP address
func (r *RealScanner) scanIP(ip string) (Device, bool) {
	// Try to ping the host
	if !r.ping(ip) {
		return Device{}, false
	}

	// Host is alive, create basic device record
	device := Device{
		IP:       ip,
		MAC:      "", // MAC discovery requires ARP, not implemented yet
		Hostname: r.lookupHostname(ip),
		Vendor:   "Unknown",
		Type:     "unknown",
	}

	return device, true
}

// ping sends ICMP ping to check if host is alive
func (r *RealScanner) ping(ip string) bool {
	var cmd *exec.Cmd

	// Platform-specific ping command
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", fmt.Sprintf("%d", r.timeoutMs), ip)
	case "darwin", "linux":
		timeoutSec := r.timeoutMs / 1000
		if timeoutSec < 1 {
			timeoutSec = 1
		}
		cmd = exec.Command("ping", "-c", "1", "-W", fmt.Sprintf("%d", timeoutSec), ip)
	default:
		return false
	}

	err := cmd.Run()
	return err == nil
}

// lookupHostname attempts reverse DNS lookup
func (r *RealScanner) lookupHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	// Remove trailing dot from FQDN
	return strings.TrimSuffix(names[0], ".")
}

// generateIPs generates all IP addresses in a CIDR range
func generateIPs(ipNet *net.IPNet) []string {
	var ips []string

	// Get network and broadcast addresses
	ip := ipNet.IP.Mask(ipNet.Mask)

	for {
		// Skip network address (first IP)
		if !ip.Equal(ipNet.IP.Mask(ipNet.Mask)) {
			ips = append(ips, ip.String())
		}

		// Increment IP
		for i := len(ip) - 1; i >= 0; i-- {
			ip[i]++
			if ip[i] != 0 {
				break
			}
		}

		// Check if we've exceeded the network range
		if !ipNet.Contains(ip) {
			break
		}
	}

	return ips
}
