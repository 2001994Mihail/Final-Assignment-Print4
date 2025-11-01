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

// Количество шагов и время прогулки
func parsePackage(data string) (int, time.Duration, error) {
	split := strings.Split(data, ",")
	if len(split) != 2 {
		err := fmt.Errorf("Неверная длинна слайса")
		return 0, 0, err
	}
	steps, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, err
	}
	if steps < 1 {
		err := fmt.Errorf("Неверное количество шагов")
		return 0, 0, err
	}
	time, err := time.ParseDuration(split[1])
	if err != nil {
		return 0, 0, err
	}
	if time < 1 {
		err := fmt.Errorf("Неверное количество времени")
		return 0, 0, err
	}
	return steps, time, nil
	// TODO: реализовать функцию
}

// Количество шагов, дистанция и калории
func DayActionInfo(data string, weight, height float64) string {
	steps, time, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps < 1 {
		return ""
	}
	distance := float64(steps) * stepLength
	distance /= mInKm
	calo, err := spentcalories.WalkingSpentCalories(steps, weight, height, time)
	if err != nil {
		return ""
	}
	string := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calo)
	return string
	// TODO: реализовать функцию
}
