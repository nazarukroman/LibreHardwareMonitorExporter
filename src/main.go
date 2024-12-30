package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	prefix := getEnv("PREFIX", "default_prefix")

	sensors, err := fetchSensors()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch sensors: %v", err), http.StatusInternalServerError)
		return
	}

	if len(sensors.Children) == 0 {
		http.Error(w, "No sensors data found", http.StatusNotFound)
		return
	}

	hostname := getHostName(sensors)
	metrics := sensorsToPrometheusMetrics(hostname, &sensors, prefix)

	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, strings.Join(metrics, "\n"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	port := getEnv("PORT", "3000")

	http.HandleFunc("/metrics", metricsHandler)

	fmt.Printf("Server is running on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
