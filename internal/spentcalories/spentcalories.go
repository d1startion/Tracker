package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных")
	}

	// шаги
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверное количество шагов")
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("неверное количество шагов - должно быть > 0")
	}
	activity := parts[1]
	durationStr := parts[2]
	if strings.Contains(durationStr, " ") {
		return 0, "", 0, fmt.Errorf("неверная продолжительность - пробел недопустим")
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверная продолжительность")
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("неверная продолжительность - должно быть > 0")
	}
	return steps, activity, duration, nil
}
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepSize := float64(height) * stepLengthCoefficient
	path := float64(steps) * stepSize
	path = path / float64(mInKm)
	return path
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	passedPath := distance(steps, height)
	hoursDuration := duration.Hours()
	return passedPath / hoursDuration
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var distanceInf, meanSpeedInf, caloriesCount float64

	switch activity {
	case "Ходьба":
		distanceInf = distance(steps, height)
		meanSpeedInf = meanSpeed(steps, height, duration)
		caloriesCount, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		distanceInf = distance(steps, height)
		meanSpeedInf = meanSpeed(steps, height, duration)
		caloriesCount, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distanceInf,
		meanSpeedInf,
		caloriesCount,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("вес и рост должны быть больше нуля")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность тренировки должна быть больше нуля")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	amount := (weight * speed * durationInMinutes) / minInH
	return amount, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("вес и рост должны быть больше нуля")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность тренировки должна быть больше нуля")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	amount := (weight * speed * durationInMinutes) / minInH
	return walkingCaloriesCoefficient * amount, nil
}
