package scheduler

import (
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/stretchr/testify/assert"
)

func TestShouldPost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p, _ := time.Parse("2021-07-18 12:55:50 +0200 SAST", "2021-07-18 12:55:50 +0200 SAST")

		tweet := &domain.Tweet{
			PostTime: p,
		}

		assert.Equal(t, true, ShouldPost(*tweet))
	})
}
