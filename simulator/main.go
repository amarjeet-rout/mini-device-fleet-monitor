package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const serverURL = "http://localhost:8080"

type heartbeatRequest struct {
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

var devices = []string{
	"device-01",
	"device-02",
	"device-03",
	"device-04",
	"device-05",
}

var (
	mu      sync.Mutex
	running = make(map[string]context.CancelFunc)
)

func main() {
	fmt.Println("Starting device simulator...")

	for _, id := range devices {
		if err := registerDevice(id); err != nil {
			fmt.Printf("failed to register %s: %v\n", id, err)
			continue
		}

		fmt.Printf("registered %s\n", id)
	}

	fmt.Println("Commands: stop <device-id>, start <device-id>, quit")

	for _, id := range devices {
		startDevice(id)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := strings.Fields(scanner.Text())

		if len(command) == 0 {
			continue
		}

		switch command[0] {
		case "stop":
			if len(command) != 2 {
				fmt.Println("usage: stop <device-id>")
				continue
			}

			stopDevice(command[1])

		case "start":
			if len(command) != 2 {
				fmt.Println("usage: start <device-id>")
				continue
			}

			startDevice(command[1])

		case "quit":
			stopAllDevices()
			fmt.Println("Simulator stopped.")
			return

		default:
			fmt.Println("unknown command")
		}
	}
}

func registerDevice(id string) error {
	body := map[string]string{
		"id":   id,
		"name": "Simulated " + id,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		serverURL+"/devices",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

func startDevice(id string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := running[id]; exists {
		fmt.Printf("%s is already running\n", id)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	running[id] = cancel

	go heartbeatLoop(ctx, id)

	fmt.Printf("%s started\n", id)
}

func stopDevice(id string) {
	mu.Lock()
	defer mu.Unlock()

	cancel, exists := running[id]
	if !exists {
		fmt.Printf("%s is not running\n", id)
		return
	}

	cancel()
	delete(running, id)

	fmt.Printf("%s stopped\n", id)
}

func stopAllDevices() {
	mu.Lock()
	defer mu.Unlock()

	for id, cancel := range running {
		cancel()
		delete(running, id)
	}
}

func heartbeatLoop(ctx context.Context, id string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	sendHeartbeat(id)

	for {
		select {
		case <-ticker.C:
			sendHeartbeat(id)

		case <-ctx.Done():
			return
		}
	}
}

func sendHeartbeat(id string) {
	request := heartbeatRequest{
		Timestamp: time.Now().UTC(),
		Status:    "healthy",
	}

	data, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("%s: failed to create heartbeat: %v\n", id, err)
		return
	}

	resp, err := http.Post(
		serverURL+"/devices/"+id+"/heartbeat",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		fmt.Printf("%s: heartbeat failed: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf(
		"%s heartbeat sent at %s\n",
		id,
		request.Timestamp.Format("15:04:05"),
	)
}
