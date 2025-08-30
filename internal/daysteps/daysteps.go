package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sc "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат данных: %q", data)
	}

	numSteps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректное количество шагов: %v", parts[0])
	}
	if numSteps <= 0 {
		return 0, 0, fmt.Errorf("некорректное количество шагов: %d", numSteps)
	}

	timeDur, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректная продолжительность: %v", parts[1])
	}
	if timeDur <= 0 {
		return 0, 0, fmt.Errorf("некорректная продолжительность: %v", parts[1])
	}

	return numSteps, timeDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepCount, timeCount, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка в данных:", err)
		return ""
	}

	distance := float64(stepCount) * stepLength / mInKm

	caloCount, err := sc.WalkingSpentCalories(stepCount, weight, height, timeCount)
	if err != nil {
		log.Println("Ошибка расчёта калорий:", err)
		return ""
	}

	// ⚠️ формат под тесты: каждая строка с новой строки + \n в конце
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		stepCount,
		distance,
		caloCount,
	)

	return result
}
