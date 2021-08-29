package webhook

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain/tweets"
	"github.com/stretchr/testify/assert"
)

func TestDetermineScheduleType(t *testing.T) {
	_ = os.Setenv("SCHEDULE_TYPE", "")
	GetSchedules()

	t.Run("Default Case", func(t *testing.T) {
		postTIme := time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local)

		expected := time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local)
		result := DetermineScheduleType(postTIme)

		assert.Equal(t, expected, result)
	})
}

func TestFixedScheduler(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.Local)
	schedules := []string{"14:30", "15:31"}
	_ = os.Setenv("SCHEDULE_TYPE", "FIXED")
	_ = os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	GetSchedules()

	t.Run("Returns next time slot", func(t *testing.T) {
		postTime := time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local)

		expected := time.Date(2021, 8, 20, 15, 31, 0, 0, time.Local)
		result := DetermineScheduleType(postTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 31, 0, 0, time.Local),
		}

		expected := time.Date(2021, 8, 21, 14, 30, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 31, 0, 0, time.Local),
		}

		expected := time.Date(2021, 9, 01, 14, 30, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestRandomMinuteScheduler(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.Local)
	seedVal = 1
	schedules := []string{"14", "15"}
	_ = os.Setenv("SCHEDULE_TYPE", "RANDOM_MINUTE")
	_ = os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	GetSchedules()

	t.Run("Returns next time slot", func(t *testing.T) {
		postTime := time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local)

		expected := time.Date(2021, 8, 20, 15, 55, 0, 0, time.Local)

		result := DetermineScheduleType(postTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 30, 0, 0, time.Local),
		}

		expected := time.Date(2021, 8, 21, 14, 55, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 30, 0, 0, time.Local),
		}

		expected := time.Date(2021, 9, 01, 14, 55, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestIntervalScheduler(t *testing.T) {
	// Ensure test falls into the switch which returns [IntervalScheduler]
	_ = os.Setenv("SCHEDULE_TYPE", "INTERVALS")
	_ = os.Setenv("INTERVALS", "2")

	t.Run("Success", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local),
		}

		expected := time.Date(2021, 8, 20, 16, 30, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next day", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 20, 23, 30, 0, 0, time.Local),
		}

		expected := time.Date(2021, 8, 21, 01, 30, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &tweets.Tweet{
			PostTime: time.Date(2021, 8, 31, 23, 30, 0, 0, time.Local),
		}

		expected := time.Date(2021, 9, 01, 01, 30, 0, 0, time.Local)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}
