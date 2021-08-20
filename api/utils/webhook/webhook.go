package webhook

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dannav/hhmmss"
)

type Schedules string

const (
	RandMin   = Schedules("RANDOM_MINUTE")
	Fixed     = Schedules("FIXED")
	Intervals = Schedules("INTERVALS")
)

func DetermineScheduleType(time time.Time) time.Time {
	typeEnv := os.Getenv("SCHEDULE_TYPE")

	switch Schedules(typeEnv) {
	case RandMin:
		return RandomMinuteScheduler(time)
	case Fixed:
		return FixedScheduler(time)
	case Intervals:
		return IntervalScheduler(time)
	default:
		return time
	}
}

func RandomMinuteScheduler(time time.Time) time.Time {

	return time
}

func FixedScheduler(inTime time.Time) time.Time {
	var result time.Time

	schedules := GetSchedules()
	schedCnt := len(schedules)
	compTime := fmt.Sprintf("%v:%v:00", inTime.Hour(), inTime.Minute())
	compDuration, _ := hhmmss.Parse(compTime)
	idx := indexOf(schedules, compDuration)
	hasNextSlot := idx+1 < schedCnt

	if hasNextSlot {
		slot := schedules[idx+1]
		result = time.Date(inTime.Year(), inTime.Month(), inTime.Day(), int(slot.Hours()), int(slot.Minutes()), inTime.Second(), inTime.Nanosecond(), inTime.Location())
	} else {
		d := inTime.Day() + 1
		slot := schedules[0]
		result = time.Date(inTime.Year(), inTime.Month(), d, int(slot.Hours()), int(slot.Minutes()), inTime.Second(), inTime.Nanosecond(), inTime.Location())

	}

	return result
}
func IntervalScheduler(time time.Time) time.Time {
	return time
}

func GetSchedules() []time.Duration {
	schedules := os.Getenv("SCHEDULES")
	slots := strings.Split(schedules, ",")
	result := make([]time.Duration, 0)

	for _, val := range slots {
		parsed, _ := hhmmss.Parse(fmt.Sprintf("%s:00", val))
		result = append(result, parsed)
	}

	return result
}

func GetInterval() int64 {
	intervals := os.Getenv("INTERVALS")

	i, err := strconv.ParseInt(intervals, 10, 64)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return i
}

func indexOf(s []time.Duration, str time.Duration) int {
	for idx, v := range s {
		if v == str {
			return idx
		}
	}

	return -1
}
