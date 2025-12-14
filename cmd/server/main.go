package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Serve static files from web/ directory
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting NetMap server on %s", addr)
	log.Printf("Serving static files from web/")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
