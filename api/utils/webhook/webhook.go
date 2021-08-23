package webhook

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Schedules string

const (
	RandMin   = Schedules("RANDOM_MINUTE")
	Fixed     = Schedules("FIXED")
	Intervals = Schedules("INTERVALS")
)

var (
	seedVal = time.Now().UnixNano()
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

func RandomMinuteScheduler(inTime time.Time) time.Time {
	schedules := GetSchedules()
	schedCnt := len(schedules)
	compTime := fmt.Sprintf("%v", inTime.Hour())

	idx := indexOf(schedules, compTime)
	hasNextSlot := idx+1 < schedCnt

	min := randomMinute(seedVal)

	if hasNextSlot {
		hour, _ := strconv.ParseInt(schedules[idx+1], 10, 64)
		return time.Date(inTime.Year(), inTime.Month(), inTime.Day(), int(hour), min, 0, 0, inTime.Location())
	} else {
		d := inTime.Day() + 1
		hour, _ := strconv.ParseInt(schedules[0], 10, 64)
		return time.Date(inTime.Year(), inTime.Month(), d, int(hour), min, 0, 0, inTime.Location())
	}
}

func FixedScheduler(inTime time.Time) time.Time {
	schedules := GetSchedules()
	schedCnt := len(schedules)
	compTime := fmt.Sprintf("%v:%v", inTime.Hour(), inTime.Minute())
	idx := indexOf(schedules, compTime)
	hasNextSlot := idx+1 < schedCnt

	if hasNextSlot {
		hour, min := splitTimeString(schedules[idx+1])

		return time.Date(inTime.Year(), inTime.Month(), inTime.Day(), hour, min, 0, 0, inTime.Location())
	} else {
		d := inTime.Day() + 1
		hour, min := splitTimeString(schedules[0])

		return time.Date(inTime.Year(), inTime.Month(), d, hour, min, 0, 0, inTime.Location())
	}
}

func IntervalScheduler(inTime time.Time) time.Time {
	interval := GetInterval()
	duration, _ := time.ParseDuration(fmt.Sprintf("%vh", interval))
	return inTime.Add(duration)
}

func GetSchedules() []string {
	schedules := os.Getenv("SCHEDULES")
	slots := strings.Split(schedules, ",")
	result := make([]string, 0)

	for _, val := range slots {
		result = append(result, val)
	}

	return result
}

func GetInterval() int64 {
	intervals := os.Getenv("INTERVALS")

	i, err := strconv.ParseInt(intervals, 10, 64)

	if err != nil {
		return 0
	}

	return i
}

func splitTimeString(slot string) (int, int) {
	hm := strings.Split(slot, ":")
	hour, _ := strconv.ParseInt(hm[0], 10, 64)
	min, _ := strconv.ParseInt(hm[1], 10, 64)

	return int(hour), int(min)
}

func randomMinute(seedVal int64) int {
	rand.Seed(seedVal)
	r := rand.Intn(59 - 1)
	return r
}

func indexOf(s []string, str string) int {
	for idx, v := range s {
		if v == str {
			return idx
		}
	}

	return -1
}
