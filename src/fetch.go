package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// fetchSensors выполняет запрос к серверу и возвращает данные о сенсорах
func fetchSensors() (Sensor, error) {
	libreHost := getEnv("LIBRE_IP", "localhost")

	url := fmt.Sprintf("http://%s/data.json", libreHost)
	resp, err := http.Get(url)
	if err != nil {
		return Sensor{}, fmt.Errorf("error making GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Sensor{}, fmt.Errorf("received non-200 HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sensor{}, fmt.Errorf("error reading response body: %w", err)
	}

	var sensors Sensor
	if err = json.Unmarshal(body, &sensors); err != nil {
		return Sensor{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return sensors, nil
}
