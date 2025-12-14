package scanner

// Device represents a discovered network device
type Device struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Vendor   string `json:"vendor"`
	Type     string `json:"type"` // router, switch, computer, printer, unknown
}

// Scanner is the interface for network scanning implementations
type Scanner interface {
	// Scan discovers devices on the network
	// Returns a channel that emits devices as they are discovered
	Scan(network string) (<-chan Device, error)
}

