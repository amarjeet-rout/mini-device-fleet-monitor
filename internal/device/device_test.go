package device

import (
	"testing"
	"time"
)

func TestDeviceStatus(t *testing.T) {
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

	device := Device{
		ID:            "device-01",
		Name:          "Lab Device 01",
		LastHeartbeat: heartbeatTime,
	}

	tests := []struct {
		name   string
		device Device
		now    time.Time
		want   string
	}{
		{
			name:   "online within timeout",
			device: device,
			now:    heartbeatTime.Add(29 * time.Second),
			want:   StatusOnline,
		},
		{
			name:   "online exactly at timeout",
			device: device,
			now:    heartbeatTime.Add(30 * time.Second),
			want:   StatusOnline,
		},
		{
			name:   "offline after timeout",
			device: device,
			now:    heartbeatTime.Add(31 * time.Second),
			want:   StatusOffline,
		},
		{
			name: "offline when no heartbeat exists",
			device: Device{
				ID:   "device-02",
				Name: "Lab Device 02",
			},
			now:  heartbeatTime,
			want: StatusOffline,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			got := tt.device.Status(tt.now)

			if got != tt.want {
				t.Errorf("Status() = %q, want %q", got, tt.want)
			}
		})
	}
}