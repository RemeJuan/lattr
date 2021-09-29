package webhook

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
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
	currentTime = time.Now().Local()
	slotsLength = 30
	timeSlots   []time.Time
)

func DetermineScheduleType(time time.Time) time.Time {
	typeEnv := os.Getenv("SCHEDULE_TYPE")

	getValidTimeSlots(time)

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
	min := randomMinute(seedVal)

	slot := timeSlots[0]
	defer removeUsedSlot(0)
	return time.Date(slot.Year(), slot.Month(), slot.Day(), slot.Hour(), min, 0, 0, inTime.Location())
}

func FixedScheduler(inTime time.Time) time.Time {
	slot := timeSlots[0]
	defer removeUsedSlot(0)
	return time.Date(slot.Year(), slot.Month(), slot.Day(), slot.Hour(), slot.Minute(), 0, 0, inTime.Location())
}

func IntervalScheduler(inTime time.Time) time.Time {
	interval := getInterval()
	duration, _ := time.ParseDuration(fmt.Sprintf("%vh", interval))
	return inTime.Add(duration)
}

func GetSchedules() {
	schedules := os.Getenv("SCHEDULES")
	slots := strings.Split(schedules, ",")
	result := make([]time.Time, 0)
	enabledDays := getScheduleDays()

	for _, val := range slots {
		hour, min := splitTimeString(val)

		for i := 0; i < slotsLength; i++ {
			t := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+i, hour, min, 0, 0, time.Local)

			if containsDay(enabledDays, t.Weekday()) {
				result = append(result, t)
			}
		}
	}

	timeSlots = result
}

func getInterval() int64 {
	intervals := os.Getenv("INTERVALS")

	i, _ := strconv.ParseInt(intervals, 10, 0)

	return i
}

// getValidTimeSlots Get time slots that are both greater than now and the most recent scheduled post
func getValidTimeSlots(lastPostTime time.Time) {
	slots := make([]time.Time, 0)

	for _, val := range timeSlots {
		if val.After(currentTime) && val.After(lastPostTime) {
			slots = append(slots, val)
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].Before(slots[j])
	})

	timeSlots = slots
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

func removeUsedSlot(i int) {
	s := timeSlots

	if i != len(s)-1 {
		s[i] = s[len(s)-1]
	}

	timeSlots = s[:len(s)-1]
}

func getScheduleDays() []time.Weekday {
	slots := make([]time.Weekday, 0)
	envDays := os.Getenv("SCHEDULE_DAYS")
	days := strings.Split(envDays, ",")

	for _, val := range days {
		switch strings.ToLower(val) {
		case "monday":
			slots = append(slots, time.Monday)
		case "tuesday":
			slots = append(slots, time.Tuesday)
		case "wednesday":
			slots = append(slots, time.Wednesday)
		case "thursday":
			slots = append(slots, time.Thursday)
		case "friday":
			slots = append(slots, time.Friday)
		case "saturday":
			slots = append(slots, time.Saturday)
		case "sunday":
			slots = append(slots, time.Sunday)
		}
	}

	return slots
}

func containsDay(arr []time.Weekday, str time.Weekday) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
