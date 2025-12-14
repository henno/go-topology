package scan

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/henno/go-topology/internal/scanner"
)

// Status represents the current state of a scan
type Status string

const (
	StatusScanning  Status = "scanning"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
	StatusFailed    Status = "failed"
)

// Scan represents a network scan instance
type Scan struct {
	ID              string            `json:"id"`
	Network         string            `json:"network"`
	CoreSwitch      string            `json:"core_switch"`
	Status          Status            `json:"status"`
	DiscoveredCount int               `json:"discovered_count"`
	Devices         []scanner.Device  `json:"devices,omitempty"`
	StartedAt       time.Time         `json:"started_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	Error           string            `json:"error,omitempty"`
	cancel          context.CancelFunc
}

// Manager handles scan lifecycle
type Manager struct {
	scanner     scanner.Scanner
	currentScan *Scan
	mu          sync.RWMutex
}

// NewManager creates a new scan manager
func NewManager(s scanner.Scanner) *Manager {
	return &Manager{
		scanner: s,
	}
}

// StartScan initiates a new network scan
func (m *Manager) StartScan(network, coreSwitch string) (*Scan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if a scan is already running
	if m.currentScan != nil && m.currentScan.Status == StatusScanning {
		return nil, fmt.Errorf("scan already in progress")
	}

	// Create new scan
	ctx, cancel := context.WithCancel(context.Background())
	scan := &Scan{
		ID:              uuid.New().String(),
		Network:         network,
		CoreSwitch:      coreSwitch,
		Status:          StatusScanning,
		DiscoveredCount: 0,
		Devices:         []scanner.Device{},
		StartedAt:       time.Now(),
		cancel:          cancel,
	}

	m.currentScan = scan

	// Start scanning in background
	go m.runScan(ctx, scan)

	return scan, nil
}

// runScan executes the scan in a goroutine
func (m *Manager) runScan(ctx context.Context, scan *Scan) {
	deviceChan, err := m.scanner.Scan(scan.Network)
	if err != nil {
		m.mu.Lock()
		scan.Status = StatusFailed
		scan.Error = err.Error()
		now := time.Now()
		scan.CompletedAt = &now
		m.mu.Unlock()
		return
	}

	// Collect devices from channel
	for {
		select {
		case <-ctx.Done():
			// Scan was cancelled
			m.mu.Lock()
			scan.Status = StatusCancelled
			now := time.Now()
			scan.CompletedAt = &now
			m.mu.Unlock()
			return

		case device, ok := <-deviceChan:
			if !ok {
				// Channel closed, scan completed
				m.mu.Lock()
				scan.Status = StatusCompleted
				now := time.Now()
				scan.CompletedAt = &now
				m.mu.Unlock()
				return
			}

			// Add device to scan results
			m.mu.Lock()
			scan.Devices = append(scan.Devices, device)
			scan.DiscoveredCount = len(scan.Devices)
			m.mu.Unlock()
		}
	}
}

// GetScan retrieves a scan by ID
func (m *Manager) GetScan(id string) (*Scan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentScan == nil || m.currentScan.ID != id {
		return nil, fmt.Errorf("scan not found")
	}

	return m.currentScan, nil
}

// GetCurrentScan retrieves the current scan
func (m *Manager) GetCurrentScan() (*Scan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentScan == nil {
		return nil, fmt.Errorf("no scan in progress")
	}

	return m.currentScan, nil
}

// CancelScan cancels a running scan
func (m *Manager) CancelScan(id string) (*Scan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentScan == nil || m.currentScan.ID != id {
		return nil, fmt.Errorf("scan not found")
	}

	if m.currentScan.Status != StatusScanning {
		return nil, fmt.Errorf("scan is not running")
	}

	// Cancel the scan
	if m.currentScan.cancel != nil {
		m.currentScan.cancel()
	}

	return m.currentScan, nil
}

