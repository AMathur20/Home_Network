package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AMathur20/Home_Network/internal/api"
	"github.com/AMathur20/Home_Network/internal/config"
	"github.com/AMathur20/Home_Network/internal/poller"
	"github.com/AMathur20/Home_Network/internal/storage"
	"github.com/AMathur20/Home_Network/internal/topology"
)

func main() {
	log.Println("Starting Home Network Monitor (HNM) v2...")

	// 1. Environment Variables & Paths
	configPath := os.Getenv("HNM_CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	dbPath := os.Getenv("HNM_DB_PATH")
	if dbPath == "" {
		dbPath = "data/hnm.db"
	}
	topoPath := "config/topology.yaml" // Typically in the same config dir
	uiPath := os.Getenv("HNM_UI_PATH")
	if uiPath == "" {
		uiPath = "ui/dist"
	}

	// Ensure data directory exists
	dataDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// 2. Load Configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Configuration loaded from %s", configPath)

	// 3. Initialize Storage (DuckDB)
	store, err := storage.NewDuckDBStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()
	log.Printf("Storage initialized at %s", dbPath)

	// 4. Load Topology
	topo, err := topology.LoadTopology(topoPath)
	if err != nil {
		log.Fatalf("Failed to load topology: %v", err)
	}
	log.Printf("Topology loaded with %d links", len(topo.Links))

	// 5. Initialize Polling Engine
	engine := poller.NewPollingEngine(cfg, store, topo)

	// 6. Watch Topology for Hot-Reload
	topology.WatchTopology(topoPath, func() {
		newTopo, err := topology.LoadTopology(topoPath)
		if err != nil {
			log.Printf("Error reloading topology: %v", err)
			return
		}
		engine.ReloadTopology(newTopo)
	})

	// 7. Start Polling Engine
	go engine.Start()

	// 8. Setup HTTP Server
	handler := api.NewAPIHandler(topoPath, store)
	
	http.HandleFunc("/api/topology", handler.GetTopology)
	http.HandleFunc("/api/metrics", handler.GetLiveMetrics)

	// Serve Static UI Files
	fs := http.FileServer(http.Dir(uiPath))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
