package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку тренировки в формате "steps,type,duration"
func parseTraining(data string) (int, string, time.Duration, error) {
	if data == "" {
		return 0, "", 0, errors.New("empty input")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid format: expected 3 parameters")
	}

	// Парсинг шагов
	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr == "" {
		return 0, "", 0, errors.New("steps cannot be empty")
	}

	// Убираем плюс если есть
	if stepsStr[0] == '+' {
		stepsStr = stepsStr[1:]
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("invalid steps format")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	// Парсинг типа тренировки
	trainingType := strings.TrimSpace(parts[1])
	if trainingType == "" {
		return 0, "", 0, errors.New("training type cannot be empty")
	}

	// Парсинг продолжительности
	durationStr := strings.TrimSpace(parts[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("duration cannot be empty")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return steps, trainingType, duration, nil
}

// distance вычисляет пройденную дистанцию в километрах
func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient
	// Дистанция в метрах
	distanceMeters := float64(steps) * stepLength
	// Дистанция в километрах
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
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

// RunningSpentCalories вычисляет количество сожженных калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле: (вес * скорость * продолжительность_в_минутах) / минут_в_часе
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories вычисляет количество сожженных калорий при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле и применяем коэффициент
	calories := (weight * speed * durationMinutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo возвращает форматированную информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	// Проверяем вес и рост
	if weight <= 0 {
		return "", errors.New("weight must be positive")
	}
	if height <= 0 {
		return "", errors.New("height must be positive")
	}

	// Вычисляем показатели
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch trainingType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	// Форматируем вывод
	info := fmt.Sprintf("Тип тренировки: %s\n", trainingType)
	info += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	info += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	info += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	info += fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

	return info, nil
}
