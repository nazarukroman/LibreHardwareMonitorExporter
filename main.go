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

//// Рекурсивная функция для преобразования в метрики
//func sensorsToPrometheusMetrics(sensors *Sensor, prefix string) string {
//	var metrics []string
//
//	// Генерация метрик для текущего сенсора
//	if value, err := strconv.ParseFloat(sensors.Value, 64); err == nil {
//		metrics = append(metrics, fmt.Sprintf(
//			"%s_value{id=\"%d\",text=\"%s\"} %f",
//			prefix, sensors.id, sensors.Text, value,
//		))
//	}
//
//	if min, err := strconv.ParseFloat(sensors.Min, 64); err == nil {
//		metrics = append(metrics, fmt.Sprintf(
//			"%s_min{id=\"%d\",text=\"%s\"} %f",
//			prefix, sensors.id, sensors.Text, min,
//		))
//	}
//
//	if max, err := strconv.ParseFloat(sensors.Max, 64); err == nil {
//		metrics = append(metrics, fmt.Sprintf(
//			"%s_max{id=\"%d\",text=\"%s\"} %f",
//			prefix, sensors.id, sensors.Text, max,
//		))
//	}
//
//	// Обработка дочерних элементов
//	for _, child := range sensors.Children {
//		metrics = append(metrics, sensorsToPrometheusMetrics(child, prefix))
//	}
//
//	fmt.Println("LOG", metrics)
//
//	return strings.Join(metrics, "\n")
//}
//
//func metricsHandler(w http.ResponseWriter, r *http.Request) {
//	// Вызываем fetchSensors
//	sensors, err := fetchSensors()
//
//	if err != nil {
//		// Обрабатываем ошибку
//		http.Error(w, fmt.Sprintf("Failed to fetch sensors: %v", err), http.StatusInternalServerError)
//		return
//	}
//
//	// Генерируем метрики
//	metrics := sensorsToPrometheusMetrics(&sensors, "libre")
//
//	// Отправляем метрики в ответ
//	w.Header().Set("Content-Type", "text/plain")
//	_, _ = w.Write([]byte(metrics))
//}
//
//func main() {
//	http.HandleFunc("/metrics", metricsHandler)
//
//	fmt.Println("Server is running on http://localhost:8080")
//
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		fmt.Println("Error starting server:", err)
//	}
//
//}

// Функция для создания SensorId
func makeSensorId(prefix, id, sensorType string) string {
	splitId := strings.Split(id, "/")
	sensorsWithoutEndingIndex := splitId[:len(splitId)-1]
	return fmt.Sprintf("%s_%s_%s", prefix, strings.Join(sensorsWithoutEndingIndex, "_"), strings.ToLower(sensorType))
}

// Функция для создания метрики
func makeMetric(sensor Sensor) string {
	return fmt.Sprintf(`{host="winnerborn-desktop",objectname="%s"} %s`, sensor.Text, sensor.Value)
}

// Функция для подготовки метрик
func prepare(sensor Sensor, prefix string) []string {
	var metrics []string

	// Проверяем, есть ли SensorId и Type
	if sensor.SensorId != nil && sensor.Type != nil {
		// Генерируем sensorId
		sensorId := makeSensorId(prefix, *sensor.SensorId, *sensor.Type)
		metric := makeMetric(sensor)

		metrics = append(metrics, sensorId, metric)
	} else if len(sensor.Children) > 0 {
		// Если есть дочерние элементы, рекурсивно вызываем prepare для каждого дочернего сенсора
		for _, child := range sensor.Children {
			metrics = append(metrics, prepare(*child, prefix)...)
		}
	}

	return metrics
}

func main() {
	sensor, err := fetchSensors()

	if err != nil {
	}

	// Подготовка метрик
	metrics := prepare(sensor, "libre")

	// Выводим метрики
	for _, metric := range metrics {
		fmt.Println(metric)
	}
}
