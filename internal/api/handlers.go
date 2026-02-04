package api

import (
	"encoding/json"
	"net/http"

	"github.com/AMathur20/Home_Network/internal/models"
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
	metrics, err := h.storage.GetLatestMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metrics == nil {
		metrics = []models.InterfaceMetric{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *APIHandler) GetMetricHistory(w http.ResponseWriter, r *http.Request) {
	device := r.URL.Query().Get("device")
	iface := r.URL.Query().Get("interface")
	if device == "" || iface == "" {
		http.Error(w, "device and interface parameters are required", http.StatusBadRequest)
		return
	}

	metrics, err := h.storage.GetMetricHistory(device, iface, 100) // Default to last 100
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
