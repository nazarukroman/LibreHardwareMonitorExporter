package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Структура Sensor с экспортируемыми полями
type Sensor struct {
	ID       int       `json:"id"`
	Text     string    `json:"Text"`
	Min      string    `json:"Min"`
	Value    string    `json:"Value"`
	Max      string    `json:"Max"`
	ImageUrl string    `json:"ImageUrl"`
	SensorId *string   `json:"SensorId"`
	Type     *string   `json:"Type"`
	Children []*Sensor `json:"Children"`
}

// Функция для получения данных с сервера
func fetchSensors() (Sensor, error) {
	libreHost := os.Getenv("LIBRE_IP")

	url := fmt.Sprintf("http://%s/data.json", libreHost)
	resp, err := http.Get(url)
	if err != nil {
		return Sensor{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sensor{}, fmt.Errorf("error reading response body: %w", err)
	}

	var sensors Sensor
	err = json.Unmarshal(body, &sensors)
	if err != nil {
		return Sensor{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return sensors, nil
}

// Функция для создания SensorId
func makeSensorId(prefix, id, text string) string {
	splitId := strings.Split(id, "/")
	sensorsWithoutEndingIndex := splitId[:len(splitId)-1]
	joinedText := strings.Split(text, " ")

	for i := range joinedText {
		joinedText[i] = strings.ToLower(joinedText[i])
	}

	return fmt.Sprintf("%s%s_%s", prefix, strings.Join(sensorsWithoutEndingIndex, "_"), strings.Join(joinedText, "_"))
}

// Функция для создания метрики
func makeMetric(hostName string, sensor Sensor) string {
	splittedText := strings.Split(sensor.Text, " ")

	for i := range splittedText {
		splittedText[i] = strings.ToLower(splittedText[i])
	}

	joinedText := strings.Join(splittedText, "_")
	valueWithoutUnit := strings.Split(sensor.Value, " ")[0]

	return fmt.Sprintf(`{host="%s",objectname="%s"} %s`, hostName, joinedText, valueWithoutUnit)
}

// Функция для подготовки метрик
func sensorsToPrometheusMetrics(hostName string, sensor *Sensor, prefix string) []string {
	var metrics []string

	if sensor.SensorId != nil && sensor.Type != nil {
		sensorId := makeSensorId(prefix, *sensor.SensorId, sensor.Text)
		metric := makeMetric(hostName, *sensor)
		metrics = append(metrics, fmt.Sprintf("%s%s", sensorId, metric))
	}

	for _, child := range sensor.Children {
		metrics = append(metrics, sensorsToPrometheusMetrics(hostName, child, prefix)...)
	}

	return metrics
}

func getHostName(sensor Sensor) string {
	return sensor.Children[0].Text
}

// Обработчик для маршрута /metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	prefix := os.Getenv("PREFIX")

	sensors, err := fetchSensors()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch sensors: %v", err), http.StatusInternalServerError)
		return
	}

	hostname := getHostName(sensors)
	metrics := sensorsToPrometheusMetrics(hostname, &sensors, prefix)
	response := strings.Join(metrics, "\n")

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(response))
}

// Главная функция
func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000" // Значение по умолчанию
	}

	http.HandleFunc("/metrics", metricsHandler)

	fmt.Printf("Server is running on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
