package webhook

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/stretchr/testify/assert"
)

func TestDetermineScheduleType(t *testing.T) {
	_ = os.Setenv("SCHEDULE_TYPE", "")

	t.Run("Default Case", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC)
		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestFixedScheduler(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.UTC)
	schedules := []string{"14:30", "15:31"}
	_ = os.Setenv("SCHEDULE_TYPE", "FIXED")
	_ = os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	t.Run("Returns next time slot", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 15, 31, 0, 0, time.UTC)
		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 31, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 14, 30, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 31, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 14, 30, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestRandomMinuteScheduler(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.UTC)
	seedVal = 1
	schedules := []string{"14", "15"}
	_ = os.Setenv("SCHEDULE_TYPE", "RANDOM_MINUTE")
	_ = os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	t.Run("Returns next time slot", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 15, 55, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 14, 55, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 14, 55, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestIntervalScheduler(t *testing.T) {
	// Ensure test falls into the switch which returns [IntervalScheduler]
	_ = os.Setenv("SCHEDULE_TYPE", "INTERVALS")
	_ = os.Setenv("INTERVALS", "2")

	t.Run("Success", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 16, 30, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next day", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 23, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 01, 30, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 23, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 01, 30, 0, 0, time.UTC)

		result := DetermineScheduleType(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestGetSchedules(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.UTC)

		first := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 14, 30, 0, 0, time.UTC)
		second := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 31, 0, 0, time.UTC)
		schedules := []time.Time{first, second}

		err := os.Setenv("SCHEDULES", strings.Join([]string{"14:30", "15:31"}, ","))
		if err != nil {
			return
		}

		expected := schedules

		result := GetSchedules()

		assert.Equal(t, expected, result)
	})
}

func TestGetInterval(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		err := os.Setenv("INTERVALS", "2")
		if err != nil {
			return
		}
		expected := int64(2)

		result := GetInterval()

		assert.Equal(t, expected, result)
	})

	t.Run("Failure", func(t *testing.T) {
		err := os.Setenv("INTERVALS", "Nan")
		if err != nil {
			return
		}
		expected := int64(0)

		result := GetInterval()

		assert.Equal(t, expected, result)
	})
}
