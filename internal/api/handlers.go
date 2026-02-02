package api

import (
	"encoding/json"
	"net/http"

	"github.com/AMathur20/Home_Network/internal/storage"
	"github.com/AMathur20/Home_Network/internal/topology"
)

type APIHandler struct {
	topoPath string
	storage  *storage.DuckDBStorage
}

func NewAPIHandler(topoPath string, s *storage.DuckDBStorage) *APIHandler {
	return &APIHandler{
		topoPath: topoPath,
		storage:  s,
	}
}

func (h *APIHandler) GetTopology(w http.ResponseWriter, r *http.Request) {
	topo, err := topology.LoadTopology(h.topoPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topo)
}

func (h *APIHandler) GetLiveMetrics(w http.ResponseWriter, r *http.Request) {
	// Simple placeholder for live metrics retrieval from DuckDB
	// In a full implementation, we'd query the latest metrics per device/interface
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Metrics retrieval not yet implemented"})
}
