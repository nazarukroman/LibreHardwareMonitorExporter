package main

import (
	"fmt"
	"regexp"
	"strings"
)

func convertToFloatString(input string) string {
	// Используем регулярное выражение для извлечения числовой части
	re := regexp.MustCompile(`^([\d,]+)`)

	// Находим соответствие с числовой частью
	match := re.FindStringSubmatch(input)

	if len(match) > 0 {
		// Заменяем запятую на точку
		result := strings.Replace(match[1], ",", ".", 1)
		return result
	}

	// Если не нашли числовое значение, возвращаем пустую строку
	return ""
}

func normalizeText(input string) string {
	// Используем регулярное выражение для удаления всех символов, кроме букв и цифр
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`) // Соответствует всем неалфавитным символам

	// Заменяем все такие символы на пробелы
	cleanedInput := re.ReplaceAllString(input, "_")

	// Разбиваем строку на слова, убираем пустые строки и приводим их к нижнему регистру
	words := strings.Fields(cleanedInput)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	// Объединяем слова в строку с подчеркиваниями
	return strings.Join(words, "_")
}

func makeSensorId(prefix, id, text string) string {
	splitId := strings.Split(id, "/")
	idPart := normalizeText(strings.Join(splitId[:len(splitId)-1], "_"))
	textPart := normalizeText(text)

	return fmt.Sprintf("%s%s_%s", prefix, idPart, textPart)
}

func makeMetric(hostName string, sensor Sensor) string {
	text := normalizeText(sensor.Text)
	value := convertToFloatString(sensor.Value)

	if value == "" {
		value = "0"
	}

	return fmt.Sprintf(`{host="%s",objectname="%s"} %s`, hostName, text, value)
}

func getHostName(sensor Sensor) string {
	if len(sensor.Children) > 0 {
		return sensor.Children[0].Text
	}
	return "unknown_host"
}
