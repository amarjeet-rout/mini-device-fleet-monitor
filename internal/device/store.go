package device

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceExists   = errors.New("device already exists")
)

type DeviceStore struct {
	mu      sync.RWMutex
	devices map[string]*Device
}

func NewDeviceStore() *DeviceStore {
	return &DeviceStore{
		devices: make(map[string]*Device),
	}
}

func (s *DeviceStore) Register(id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.devices[id]; exists {
		return ErrDeviceExists
	}

	s.devices[id] = &Device{
		ID:   id,
		Name: name,
	}

	return nil
}

func (s *DeviceStore) Get(id string) (Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, exists := s.devices[id]
	if !exists {
		return Device{}, ErrDeviceNotFound
	}

	return *device, nil
}

func (s *DeviceStore) List() []Device {
	s.mu.RLock()
	defer s.mu.RUnlock()

	devices := make([]Device, 0, len(s.devices))

	for _, device := range s.devices {
		devices = append(devices, *device)
	}

	return devices
}

func (s *DeviceStore) Heartbeat(id string, timestamp time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	device, exists := s.devices[id]
	if !exists {
		return ErrDeviceNotFound
	}

	device.LastHeartbeat = timestamp

	return nil
}