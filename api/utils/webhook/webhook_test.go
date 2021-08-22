package webhook

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/stretchr/testify/assert"
)

func TestFixedScheduler(t *testing.T) {
	schedules := []string{"14:30", "15:31"}
	err := os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("Returns next time slot", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 15, 31, 0, 0, time.UTC)
		result := FixedScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 31, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 14, 30, 0, 0, time.UTC)

		result := FixedScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 31, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 14, 30, 0, 0, time.UTC)

		result := FixedScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestRandomMinuteScheduler(t *testing.T) {
	seedVal = 1
	schedules := []string{"14", "15"}
	err := os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("Returns next time slot", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 15, 55, 0, 0, time.UTC)

		result := RandomMinuteScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next days date", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 15, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 14, 55, 0, 0, time.UTC)

		result := RandomMinuteScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 15, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 14, 55, 0, 0, time.UTC)

		result := RandomMinuteScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestIntervalScheduler(t *testing.T) {
	err := os.Setenv("INTERVALS", "2")
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("Success", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 14, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 20, 16, 30, 0, 0, time.UTC)

		result := IntervalScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Returns next day", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 20, 23, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 8, 21, 01, 30, 0, 0, time.UTC)

		result := IntervalScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})

	t.Run("Last day of month", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: time.Date(2021, 8, 31, 23, 30, 0, 0, time.UTC),
		}

		expected := time.Date(2021, 9, 01, 01, 30, 0, 0, time.UTC)

		result := IntervalScheduler(tweet.PostTime)

		assert.Equal(t, expected, result)
	})
}

func TestGetSchedules(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		schedules := []string{"14:30", "15:31"}
		err := os.Setenv("SCHEDULES", strings.Join(schedules, ","))
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
