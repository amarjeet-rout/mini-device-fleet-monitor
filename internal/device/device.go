package device

import "time"

const (
	StatusOnline  = "ONLINE"
	StatusOffline = "OFFLINE"
)

const heartbeatTimeout = 30 * time.Second

type Device struct {
	ID              string
	Name            string
	LastHeartbeat   time.Time
	HeartbeatStatus string
}

func (d Device) Status(now time.Time) string {
	if d.LastHeartbeat.IsZero() {
		return StatusOffline
	}

	if now.Sub(d.LastHeartbeat) <= heartbeatTimeout {
		return StatusOnline
	}

	return StatusOffline
}