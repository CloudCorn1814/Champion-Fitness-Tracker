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

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid string format: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("steps parsing error: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("positive steps value expected: %v", err)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("training duration parsing error: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("positive duration value expected: %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("data parsing error:", err)
		return ""
	}

	if steps <= 0 {
		log.Println("incorrect steps value:", steps)
		return ""
	}

	distanceM := stepLength * float64(steps)

	distanceKm := distanceM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("calorie calculation error:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

}
