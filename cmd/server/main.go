package main

import (
	"log"
	"net/http"

	"mini-device-fleet-monitor/internal/device"
	fleethttp "mini-device-fleet-monitor/internal/http"
)

func main() {
	store := device.NewDeviceStore()
	handler := fleethttp.NewHandler(store)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("server listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
