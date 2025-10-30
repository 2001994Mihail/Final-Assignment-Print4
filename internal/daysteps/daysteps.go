package daysteps

import (
	"fmt"
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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: expected 'steps,duration'")
	}

	// Парсим шаги
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps format: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("steps must be positive")
	}

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[1])
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
		return ""
	}

	// Вычисляем дистанцию
	distance := (float64(steps) * stepLength) / mInKm

	// Вычисляем калории
	kcal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	// Форматируем результат
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, kcal)
}
