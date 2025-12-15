package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/henno/go-topology/internal/scan"
	"github.com/henno/go-topology/internal/scanner"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Scanner struct {
		TimeoutMs int `json:"timeout_ms"`
		Workers   int `json:"workers"`
	} `json:"scanner"`
}

func loadConfig() (*Config, error) {
	configPath := os.Getenv("NETMAP_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func isMockMode() bool {
	mockEnv := os.Getenv("NETMAP_MOCK")
	return mockEnv == "true" || mockEnv == "1"
}

func mockModeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject mock mode indicator into HTML responses
		if isMockMode() && strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/" {
			w.Header().Set("X-Mock-Mode", "true")
		}
		next.ServeHTTP(w, r)
	})
}

// handleStartScan handles POST /api/scans
func handleStartScan(w http.ResponseWriter, r *http.Request, manager *scan.Manager) {
	var req struct {
		Network    string `json:"network"`
		CoreSwitch string `json:"core_switch"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Network == "" {
		http.Error(w, "network is required", http.StatusBadRequest)
		return
	}

	s, err := manager.StartScan(req.Network, req.CoreSwitch)
	if err != nil {
		if err.Error() == "scan already in progress" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// handleGetCurrentScan handles GET /api/scans/current
func handleGetCurrentScan(w http.ResponseWriter, r *http.Request, manager *scan.Manager) {
	s, err := manager.GetCurrentScan()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// handleGetScan handles GET /api/scans/{id}
func handleGetScan(w http.ResponseWriter, r *http.Request, manager *scan.Manager, scanID string) {
	s, err := manager.GetScan(scanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// handleCancelScan handles DELETE /api/scans/{id}
func handleCancelScan(w http.ResponseWriter, r *http.Request, manager *scan.Manager, scanID string) {
	s, err := manager.CancelScan(scanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	mockMode := isMockMode()
	if mockMode {
		log.Printf("🧪 Running in MOCK MODE")
	}

	// Initialize scanner
	var s scanner.Scanner
	if mockMode {
		s = scanner.NewMockScanner()
	} else {
		s = scanner.NewRealScanner(config.Scanner.Workers, config.Scanner.TimeoutMs)
		log.Printf("🔍 Running in PRODUCTION MODE with %d workers", config.Scanner.Workers)
	}

	// Initialize scan manager
	scanManager := scan.NewManager(s)

	// Create a custom router to handle API and static files
	mux := http.NewServeMux()

	// API endpoint to check mock mode
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := map[string]interface{}{
			"mock_mode": mockMode,
		}
		json.NewEncoder(w).Encode(status)
	})

	// Scan API endpoints
	mux.HandleFunc("/api/scans/current", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetCurrentScan(w, r, scanManager)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle POST /api/scans (exact match without trailing slash)
	mux.HandleFunc("/api/scans", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleStartScan(w, r, scanManager)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Handle GET/DELETE /api/scans/{id}
	mux.HandleFunc("/api/scans/", func(w http.ResponseWriter, r *http.Request) {
		// Extract scan ID from path
		scanID := strings.TrimPrefix(r.URL.Path, "/api/scans/")
		if scanID == "" || scanID == "current" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetScan(w, r, scanManager, scanID)
		case http.MethodDelete:
			handleCancelScan(w, r, scanManager, scanID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve HTML pages with clean URLs
	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/scan.html")
	})
	mux.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/devices.html")
	})

	// Serve static files from web/ directory
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", mockModeHandler(fs))

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting NetMap server on %s", addr)
	log.Printf("Serving static files from web/")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
