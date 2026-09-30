package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mini-device-fleet-monitor/internal/device"
)

func TestRegisterDevice(t *testing.T) {
	store := device.NewDeviceStore()
	handler := NewHandler(store)

	body := `{"id":"device-01","name":"Lab Device 01"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusCreated)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["id"] != "device-01" {
		t.Errorf("response id = %q, want %q", response["id"], "device-01")
	}

	if response["name"] != "Lab Device 01" {
		t.Errorf("response name = %q, want %q", response["name"], "Lab Device 01")
	}
}

func TestRegisterDuplicateDevice(t *testing.T) {
	store := device.NewDeviceStore()
	handler := NewHandler(store)

	body := `{"id":"device-01","name":"Lab Device 01"}`

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(body),
	)

	rec1 := httptest.NewRecorder()

	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusCreated {
		t.Fatalf("first registration status = %d, want %d", rec1.Code, http.StatusCreated)
	}

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(body),
	)

	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("duplicate registration status = %d, want %d", rec2.Code, http.StatusConflict)
	}
}

func TestRegisterDeviceInvalidJSON(t *testing.T) {
	store := device.NewDeviceStore()
	handler := NewHandler(store)

	req := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(`invalid json`),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestListDevices(t *testing.T) {
	store := device.NewDeviceStore()
	handler := NewHandler(store)

	if err := store.Register("device-01", "Lab Device 01"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := store.Register("device-02", "Lab Device 02"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/devices",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var devices []device.Device

	if err := json.NewDecoder(rec.Body).Decode(&devices); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(devices) != 2 {
		t.Errorf("number of devices = %d, want 2", len(devices))
	}
}
