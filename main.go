package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Sensor struct {
	id       int
	Text     string
	Min      string
	Value    string
	Max      string
	ImageUrl string
	SensorId *string
	Type     *string
	Children []*Sensor
}

func fetchSensors() (Sensor, error) {
	resp, err := http.Get("http://192.168.1.149:8085/data.json")

	if err != nil {
		return Sensor{}, fmt.Errorf("error making GET request: %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Sensor{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Распарсить JSON в структуру Sensor
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
func makeMetric(sensor Sensor) string {
	splittedText := strings.Split(sensor.Text, " ")

	// Преобразуем слова в нижний регистр
	for i := range splittedText {
		splittedText[i] = strings.ToLower(splittedText[i])
	}

	// Объединяем слова через символ "_"
	joinedText := strings.Join(splittedText, "_")
	valueWithoutUnit := strings.Split(sensor.Value, " ")[0]

	// Формируем метрику
	return fmt.Sprintf(`{host="winnerborn-desktop",objectname="%s"} %s`, joinedText, valueWithoutUnit)
}

// Функция для подготовки метрик
func sensorsToPrometheusMetrics(sensor Sensor, prefix string) []string {
	var metrics []string

	// Проверяем, есть ли SensorId и Type
	if sensor.SensorId != nil && sensor.Type != nil {
		// Генерируем sensorId
		sensorId := makeSensorId(prefix, *sensor.SensorId, sensor.Text)
		metric := makeMetric(sensor)

		metrics = append(metrics, sensorId, metric)
	} else if len(sensor.Children) > 0 {
		// Если есть дочерние элементы, рекурсивно вызываем prepare для каждого дочернего сенсора
		for _, child := range sensor.Children {
			metrics = append(metrics, sensorsToPrometheusMetrics(*child, prefix)...)
		}
	}

	return metrics
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Вызываем fetchSensors
	sensors, err := fetchSensors()

	if err != nil {
		// Обрабатываем ошибку
		http.Error(w, fmt.Sprintf("Failed to fetch sensors: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерируем метрики
	metrics := sensorsToPrometheusMetrics(&sensors, "libre")

	// Отправляем метрики в ответ
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(metrics))
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)

	fmt.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
