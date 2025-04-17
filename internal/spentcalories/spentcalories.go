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
		return 0, "", 0, fmt.Errorf("invalid string format: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("steps parsing error: %v", err)
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("positive steps value expected: %v", err)
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("activity duration parsing error: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("positive duration value expected: %v", duration)
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
		return "", fmt.Errorf("data parsing error: %v", err)
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	switch activityType {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			duration.Hours(), dist, speed, calories), nil

	case "Бег":
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
		return 0, fmt.Errorf("invalid number of steps provided")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight value provided")
	}

	if height <= 0 {
		return 0, fmt.Errorf("invalid height value provided")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid activity duration value provided")
	}

	spentCalories := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid number of steps provided")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight value provided")
	}

	if height <= 0 {
		return 0, fmt.Errorf("invalid height value provided")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid activity duration value provided")
	}

	spentCalories := walkingCaloriesCoefficient * ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH)

	return spentCalories, nil
}
