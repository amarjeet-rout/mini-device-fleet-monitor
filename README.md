# Mini Device Fleet Monitor

A small Go application that monitors the health of simulated devices using periodic heartbeat messages.

The system provides HTTP APIs to register devices, receive heartbeats, view device status, and view an overall fleet summary.

Device connectivity is determined automatically:

- `ONLINE` — heartbeat received within the last 30 seconds
- `OFFLINE` — no heartbeat received for more than 30 seconds## Architecture

The application is built using Go and the standard library.

### Components

- **HTTP API** — exposes device registration, heartbeat, device details, device listing, and fleet summary endpoints.
- **Device model** — contains device information and calculates `ONLINE`/`OFFLINE` status from the latest heartbeat.
- **Device store** — an in-memory, concurrency-safe store protected by `sync.RWMutex`.
- **Simulator** — creates five simulated devices and sends heartbeats every 5 seconds. Individual devices can be stopped and restarted from the simulator terminal.

### Status calculation

Connectivity status is not stored as a permanent `ONLINE` or `OFFLINE` value.

Instead, it is calculated from the latest heartbeat:

```text
heartbeat age <= 30 seconds  → ONLINE
heartbeat age > 30 seconds   → OFFLINE
no heartbeat                 → OFFLINE

## Prerequisites

- Go 1.20 or later
- Git
- PowerShell or another terminal

## Build and Run

Clone the repository and move into the project directory:

```bash
git clone <repository-url>
cd mini-device-fleet-monitor

## Simulator

The simulator creates five devices:

- `device-01`
- `device-02`
- `device-03`
- `device-04`
- `device-05`

Each device sends a heartbeat every 5 seconds.

Start the server first:

```bash
go run ./cmd/server

Then, in another terminal, start the simulator:

go run ./simulator

The simulator supports these commands:

stop <device-id>
start <device-id>
quit

For example:

stop device-03

stops heartbeats for device-03. After more than 30 seconds without a heartbeat, the API reports the device as OFFLINE.

To resume heartbeats:

start device-03
## API

### Register a device

```http
POST /devices
Content-Type: application/json

{
  "id": "device-01",
  "name": "Lab Device 01"
}
Returns 201 Created
Send a heartbeat
POST /devices/{id}/heartbeat
Content-Type: application/json

Request:

{
  "timestamp": "2026-09-30T12:00:00Z",
  "status": "healthy"
}

Returns 200 OK.

List all devices
GET /devices

Example response:

[
  {
    "id": "device-01",
    "name": "Lab Device 01",
    "status": "ONLINE",
    "last_heartbeat": "2026-09-30T12:00:00Z"
  }
]
Get a device
GET /devices/{id}

Returns the device's current status and latest heartbeat.

Get fleet summary
GET /summary

Example response:

{
  "total": 5,
  "online": 4,
  "offline": 1
}

## Tests

Run all automated tests with:

```bash
go test ./...

The test suite covers:

Device registration
Duplicate device registration
Heartbeat handling
Unknown device handling
Device status calculation
No-heartbeat status
The 30-second timeout boundary
Fleet summary
HTTP API behavior

The timeout behavior is tested without waiting 30 seconds in real time. Fixed timestamps are used to verify:

29 seconds → ONLINE
30 seconds → ONLINE
31 seconds → OFFLINE


## Assumptions

- Device data is stored in memory.
- Device IDs are unique.
- A heartbeat timestamp is treated as the device's latest heartbeat time.
- A device is considered `ONLINE` when its latest heartbeat is no more than 30 seconds old.
- A device with no heartbeat is considered `OFFLINE`.
- The application currently runs as a single server instance.

## Limitations

- Device data is lost when the server restarts.
- There is no persistent database.
- The heartbeat `status` field is validated but connectivity status is determined from heartbeat age.
- Heartbeat metrics are accepted by the API but are not currently stored or exposed.
- There is no authentication or authorization.
- The API does not currently provide pagination or advanced filtering.

## Next-Day Improvements

If this project were extended beyond the assignment, the following improvements would be considered:

- Persist devices and heartbeats in a database.
- Add structured logging.
- Add request validation and clearer API error responses.
- Add authentication and authorization.
- Add metrics and monitoring.
- Add API documentation using OpenAPI.
- Add configurable server address and heartbeat timeout.
- Add more simulator controls and configurable device counts.
- Add graceful server shutdown.
- Add a web dashboard for visualizing device status.