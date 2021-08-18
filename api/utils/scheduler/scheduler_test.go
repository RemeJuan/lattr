package scheduler

import (
	"testing"

	"github.com/RemeJuan/lattr/domain"
	"github.com/stretchr/testify/assert"
)

func TestShouldPost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: "2021-07-18 12:55:50 +0200",
		}

		assert.Equal(t, true, ShouldPost(*tweet))
	})

	t.Run("Error", func(t *testing.T) {
		tweet := &domain.Tweet{
			PostTime: "2021-07-18 12:55:50 +0200 SAST",
		}

		assert.Equal(t, false, ShouldPost(*tweet))
	})
}
