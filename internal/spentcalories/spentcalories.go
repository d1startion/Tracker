package spentcalories

import (
	"fmt"
	"log"
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
	trainingSlice := strings.Split(data, ",")
	if len(trainingSlice) != 3 {
		return 0, "", 0, strconv.ErrSyntax
	}
	trainingSteps, err := strconv.Atoi(trainingSlice[0])
	if err != nil {
		return 0, "", 0, err
	}
	timeTraining, err := time.ParseDuration(trainingSlice[1])
	if err != nil {
		return 0, "", 0, err
	}
	return trainingSteps, trainingSlice[1], timeTraining, nil
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
	trainingSteps, typeOfActivity, timeTraining, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var distanceInf, meanSpeedInf, caloriesCount float64
	switch typeOfActivity {
	case "Ходьба":
		distanceInf = distance(trainingSteps, height)
		meanSpeedInf = meanSpeed(trainingSteps, height, timeTraining)
		caloriesCount, err = WalkingSpentCalories(trainingSteps, weight, height, timeTraining)
		if err != nil {
			return "", err
		}
	case "Бег":
		distanceInf = distance(trainingSteps, height)
		meanSpeedInf = meanSpeed(trainingSteps, height, timeTraining)
		caloriesCount, err = RunningSpentCalories(trainingSteps, weight, height, timeTraining)
		if err != nil {
			return "", err
		}
	}

	result := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f",
		typeOfActivity,
		timeTraining.Hours(),
		distanceInf,
		meanSpeedInf,
		caloriesCount)
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
