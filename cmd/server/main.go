package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	mockMode := isMockMode()
	if mockMode {
		log.Printf("🧪 Running in MOCK MODE")
	}

	// Serve static files from web/ directory
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", mockModeHandler(fs))

	// API endpoint to check mock mode
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := map[string]interface{}{
			"mock_mode": mockMode,
		}
		json.NewEncoder(w).Encode(status)
	})

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting NetMap server on %s", addr)
	log.Printf("Serving static files from web/")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
