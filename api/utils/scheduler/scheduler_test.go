package scheduler

import (
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain/tweets"
	"github.com/stretchr/testify/assert"
)

const layout = "2021-07-18 12:55:50 +0200 SAST"

func TestShouldPost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p, _ := time.Parse(layout, "2021-07-18 12:55:50 +0200 SAST")

		tweet := &tweets.Tweet{
			PostTime: p,
		}

		assert.Equal(t, true, ShouldPost(*tweet))
	})
}
