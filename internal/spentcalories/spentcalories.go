package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Коэффициент длины шага (ИСПРАВЛЕНО!)
	stepLengthCoefficient = 0.45
	// Количество метров в одном километре
	mInKm = 1000
	// Количество минут в часе
	minInH = 60
	// Корректирующий коэффициент для ходьбы (ИСПРАВЛЕНО!)
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
		return 0, "", 0, fmt.Errorf("invalid steps format: %v", err)
	}

	// Вид активности
	activity := parts[1]

	// Парсим продолжительность - БЕЗ обрезки пробелов!
	durationStr := parts[2]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %v", err)
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

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return (float64(steps) * stepLength) / mInKm
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	hours := duration.Hours()

	if hours == 0 {
		return 0
	}

	return dist / hours
}

// RunningSpentCalories вычисляет калории для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет калории для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo возвращает информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var errCal error

	switch activity {
	case "Бег":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
		if errCal != nil {
			return "", errCal
		}
	case "Ходьба":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
		if errCal != nil {
			return "", errCal
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// УБИРАЕМ лишний перенос строки в конце!
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), dist, speed, calories)

	return result, nil
}
