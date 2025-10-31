package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient      = 0.45
	mInKm                      = 1000
	minInH                     = 60
	walkingCaloriesCoefficient = 0.5
)

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {
	if data == "" {
		return 0, "", 0, fmt.Errorf("empty input")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid format: expected 'steps,activity,duration'")
	}

	// Парсим шаги - БЕЗ обрезки пробелов!
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps format: %w", err)
	}

	// Вид активности - БЕЗ обрезки пробелов!
	activity := parts[1]

	// Парсим продолжительность - БЕЗ обрезки пробелов!
	durationStr := parts[2]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}

	// Проверяем валидность данных
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps must be positive")
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be positive")
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в км
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

// RunningSpentCalories вычисляет калории для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	// Формула для бега: (вес * скорость * время_в_минутах) / 60
	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет калории для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	// Формула для ходьбы: ((вес * скорость * время_в_минутах) / 60) * коэффициент
	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo анализирует тренировку и возвращает результат
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64

	// Обрабатываем разные варианты написания активности
	switch strings.ToLower(activity) {
	case "бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("unknown activity: %s", activity)
	}

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// УДАЛЕН завершающий символ новой строки
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), dist, speed, calories)

	return result, nil
}
