package main

import (
	"fmt"
	"regexp"
)

func sensorsToPrometheusMetrics(hostName string, root *Sensor, prefix string) []string {
	var metrics []string
	stack := []*Sensor{root}

	re := regexp.MustCompile(`^[a-zA-Z0-9/-]+$`)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.SensorId != nil && current.Type != nil && re.MatchString(*current.SensorId) {
			sensorId := makeSensorId(prefix, *current.SensorId, current.Text)
			metric := makeMetric(hostName, *current)
			metrics = append(metrics, fmt.Sprintf("%s%s", sensorId, metric))
		}

		stack = append(stack, current.Children...)
	}

	return metrics
}
