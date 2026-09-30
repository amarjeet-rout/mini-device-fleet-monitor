package device

import (
	"testing"
	"time"
)

func TestDeviceStoreRegisterAndGet(t *testing.T) {
	store := NewDeviceStore()

	err := store.Register("device-01", "Lab Device 01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, err := store.Get("device-01")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.ID != "device-01" {
		t.Errorf("ID = %q, want %q", got.ID, "device-01")
	}

	if got.Name != "Lab Device 01" {
		t.Errorf("Name = %q, want %q", got.Name, "Lab Device 01")
	}
}

func TestDeviceStoreDuplicateRegister(t *testing.T) {
	store := NewDeviceStore()

	err := store.Register("device-01", "Lab Device 01")
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	err = store.Register("device-01", "Another Device")
	if err != ErrDeviceExists {
		t.Errorf("second Register() error = %v, want %v", err, ErrDeviceExists)
	}
}

func TestDeviceStoreGetMissingDevice(t *testing.T) {
	store := NewDeviceStore()

	_, err := store.Get("device-404")

	if err != ErrDeviceNotFound {
		t.Errorf("Get() error = %v, want %v", err, ErrDeviceNotFound)
	}
}

func TestDeviceStoreHeartbeat(t *testing.T) {
	store := NewDeviceStore()

	err := store.Register("device-01", "Lab Device 01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	heartbeatTime := time.Date(
		2026,
		9,
		30,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	err = store.Heartbeat("device-01", heartbeatTime)
	if err != nil {
		t.Fatalf("Heartbeat() error = %v", err)
	}

	got, err := store.Get("device-01")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !got.LastHeartbeat.Equal(heartbeatTime) {
		t.Errorf("LastHeartbeat = %v, want %v", got.LastHeartbeat, heartbeatTime)
	}
}

func TestDeviceStoreList(t *testing.T) {
	store := NewDeviceStore()

	if err := store.Register("device-01", "Lab Device 01"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := store.Register("device-02", "Lab Device 02"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	devices := store.List()

	if len(devices) != 2 {
		t.Fatalf("List() returned %d devices, want 2", len(devices))
	}
}