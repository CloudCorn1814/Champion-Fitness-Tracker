package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат строки: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге количества шагов: %v", err)
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге длительности активности: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть положительной, получено: %v", duration)
	}
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	distance := ((height * stepLengthCoefficient) * float64(steps)) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	mediumSpeed := distance(steps, height) / duration.Hours()

	return mediumSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга: %v", err)
	}

	switch activityType {
	case "Ходьба":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			duration.Hours(), dist, speed, calories), nil

	case "Бег":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			duration.Hours(), dist, speed, calories), nil

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("ошибка при передаче количества шагов")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("ошибка при передаче веса")
	}

	if height <= 0 {
		return 0, fmt.Errorf("ошибка при передаче роста")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("ошибка при передаче длительности активности")
	}

	spentCalories := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("ошибка при передаче количества шагов")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("ошибка при передаче веса")
	}

	if height <= 0 {
		return 0, fmt.Errorf("ошибка при передаче роста")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("ошибка при передаче длительности активности")
	}

	spentCalories := walkingCaloriesCoefficient * ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH)

	return spentCalories, nil
}
