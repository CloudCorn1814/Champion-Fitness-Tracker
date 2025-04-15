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
		return 0, 0, fmt.Errorf("неверный формат строки: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if steps <= 0 || err != nil {
		return 0, 0, fmt.Errorf("ошибка при парсинге количества шагов: %v", err)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при парсинге длительности ходьбы: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("длительность должна быть положительной: %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceM := stepLength * float64(steps)

	distanceKm := distanceM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка при расчёте калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

}
