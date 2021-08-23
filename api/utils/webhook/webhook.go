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
	seedVal     = time.Now().UnixNano()
	currentTime = time.Now()
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
	slots := getValidTimeSlots(schedules, inTime)
	hasNextSlot := len(slots) > 0

	min := randomMinute(seedVal)

	if hasNextSlot {
		slot := slots[0]
		return time.Date(inTime.Year(), inTime.Month(), inTime.Day(), slot.Hour(), min, 0, 0, inTime.Location())
	} else {
		d := inTime.Day() + 1
		slot := schedules[0]
		return time.Date(inTime.Year(), inTime.Month(), d, slot.Hour(), min, 0, 0, inTime.Location())
	}
}

func FixedScheduler(inTime time.Time) time.Time {
	schedules := GetSchedules()
	slots := getValidTimeSlots(schedules, inTime)
	hasNextSlot := len(slots) > 0

	if hasNextSlot {
		slot := slots[0]
		return time.Date(inTime.Year(), inTime.Month(), inTime.Day(), slot.Hour(), slot.Minute(), 0, 0, inTime.Location())
	} else {
		d := inTime.Day() + 1
		slot := schedules[0]

		return time.Date(inTime.Year(), inTime.Month(), d, slot.Hour(), slot.Minute(), 0, 0, inTime.Location())
	}
}

func IntervalScheduler(inTime time.Time) time.Time {
	interval := GetInterval()
	duration, _ := time.ParseDuration(fmt.Sprintf("%vh", interval))
	return inTime.Add(duration)
}

func GetSchedules() []time.Time {
	schedules := os.Getenv("SCHEDULES")
	slots := strings.Split(schedules, ",")
	result := make([]time.Time, 0)

	for _, val := range slots {
		hour, min := splitTimeString(val)
		t := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hour, min, 0, 0, time.UTC)
		result = append(result, t)
	}

	return result
}

func GetInterval() int64 {
	intervals := os.Getenv("INTERVALS")

	i, err := strconv.ParseInt(intervals, 10, 0)

	if err != nil {
		return 0
	}

	return i
}

// Get time slots that are both greater than now and the most recent scheduled post
func getValidTimeSlots(scheduleSlots []time.Time, lastPostTime time.Time) []time.Time {
	slots := make([]time.Time, 0)

	for _, val := range scheduleSlots {
		if val.After(currentTime) && val.After(lastPostTime) {
			slots = append(slots, val)
		}
	}

	return slots
}

func splitTimeString(slot string) (int, int) {
	var minute int

	hm := strings.Split(strings.TrimSpace(slot), ":")
	hour, _ := strconv.ParseInt(hm[0], 10, 0)
	if len(hm) == 1 {
		minute = 0
	} else {
		min, _ := strconv.ParseInt(hm[1], 10, 0)
		minute = int(min)
	}

	return int(hour), minute
}

func randomMinute(seedVal int64) int {
	rand.Seed(seedVal)
	r := rand.Intn(59 - 1)
	return r
}
