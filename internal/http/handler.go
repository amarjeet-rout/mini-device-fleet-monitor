package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"mini-device-fleet-monitor/internal/device"
)

type Handler struct {
	store *device.DeviceStore
}

func NewHandler(store *device.DeviceStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	switch {
	case r.Method == http.MethodPost && path == "devices":
		h.registerDevice(w, r)

	case r.Method == http.MethodGet && path == "devices":
		h.listDevices(w, r)

	case r.Method == http.MethodGet && path == "summary":
		h.summary(w, r)

	case r.Method == http.MethodGet && strings.HasPrefix(path, "devices/"):
		h.getDevice(w, r, path)

	case r.Method == http.MethodPost &&
		strings.HasPrefix(path, "devices/") &&
		strings.HasSuffix(path, "/heartbeat"):
		h.heartbeat(w, r, path)

	default:
		http.NotFound(w, r)
	}
}

type registerDeviceRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type heartbeatRequest struct {
	Timestamp time.Time          `json:"timestamp"`
	Status    string             `json:"status"`
	Metrics   map[string]float64 `json:"metrics,omitempty"`
}

type deviceResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

type summaryResponse struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

func (h *Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	var req registerDeviceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Name == "" {
		http.Error(w, "id and name are required", http.StatusBadRequest)
		return
	}

	if err := h.store.Register(req.ID, req.Name); err != nil {
		if err == device.ErrDeviceExists {
			http.Error(w, "device already exists", http.StatusConflict)
			return
		}

		http.Error(w, "failed to register device", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":      req.ID,
		"name":    req.Name,
		"message": "device registered",
	})
}

func (h *Handler) listDevices(w http.ResponseWriter, r *http.Request) {
	devices := h.store.List()
	now := time.Now()

	response := make([]deviceResponse, 0, len(devices))

	for _, d := range devices {
		response = append(response, deviceResponse{
			ID:            d.ID,
			Name:          d.Name,
			Status:        d.Status(now),
			LastHeartbeat: d.LastHeartbeat,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) heartbeat(w http.ResponseWriter, r *http.Request, path string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 3 || parts[0] != "devices" || parts[2] != "heartbeat" {
		http.NotFound(w, r)
		return
	}

	deviceID := parts[1]

	var req heartbeatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Timestamp.IsZero() {
		http.Error(w, "timestamp is required", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if err := h.store.Heartbeat(deviceID, req.Timestamp); err != nil {
		if err == device.ErrDeviceNotFound {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "heartbeat received",
		"id":      deviceID,
	})
}

func (h *Handler) getDevice(w http.ResponseWriter, r *http.Request, path string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 2 || parts[0] != "devices" || parts[1] == "" {
		http.NotFound(w, r)
		return
	}

	deviceID := parts[1]

	d, err := h.store.Get(deviceID)
	if err != nil {
		if err == device.ErrDeviceNotFound {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get device", http.StatusInternalServerError)
		return
	}

	response := deviceResponse{
		ID:            d.ID,
		Name:          d.Name,
		Status:        d.Status(time.Now()),
		LastHeartbeat: d.LastHeartbeat,
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	devices := h.store.List()
	now := time.Now()

	response := summaryResponse{
		Total: len(devices),
	}

	for _, d := range devices {
		if d.Status(now) == device.StatusOnline {
			response.Online++
		} else {
			response.Offline++
		}
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}
