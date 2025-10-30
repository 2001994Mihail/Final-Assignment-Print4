package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку с данными о шагах и продолжительности
func parsePackage(data string) (int, time.Duration, error) {
	if data == "" {
		return 0, 0, fmt.Errorf("empty input")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: expected 'steps,duration'")
	}

	// Парсим шаги - БЕЗ обрезки пробелов!
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps format: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("steps must be positive")
	}

	// Парсим продолжительность - БЕЗ обрезки пробелов!
	durationStr := parts[1]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration must be positive")
	}

	return steps, duration, nil
}

// DayActionInfo возвращает информацию о дневной активности
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		// ДОБАВЛЯЕМ ВЫВОД ОШИБКИ В ЛОГ
		log.Println(err)
		return ""
	}

	// Вычисляем дистанцию
	distance := (float64(steps) * stepLength) / mInKm

	// Вычисляем калории
	kcal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		// ДОБАВЛЯЕМ ВЫВОД ОШИБКИ В ЛОГ
		log.Println(err)
		return ""
	}

	// УБИРАЕМ лишний перенос строки в конце!
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distance, kcal)
}
