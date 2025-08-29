package daysteps

import (
	"fmt"
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
		return 0, 0, strconv.ErrSyntax
	}
	numSteps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if numSteps <= 0 {
		return 0, 0, err
	}
	timeDur, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return numSteps, timeDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepCount, timeCount, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if stepCount <= 0 {
		return ""
	}
	distance := float64(stepCount) * stepLength
	distance = distance / mInKm
	ciloCount, err := sc.WalkingSpentCalories(stepCount, weight, height, timeCount)
	if err != nil {
		return ""
	}
	result := fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.`, stepCount, distance, ciloCount)

	return result
}
