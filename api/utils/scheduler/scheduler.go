package scheduler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RemeJuan/lattr/domain/tweets"
	"github.com/RemeJuan/lattr/utils/twitter"
	"github.com/RemeJuan/lattr/utils/webhook"
	"github.com/go-co-op/gocron"
)

func Scheduler() {
	var schedule string
	s := gocron.NewScheduler(time.Local)
	cr := os.Getenv("CRON_SCHEDULE")

	if len(cr) == 0 {
		schedule = "*/5 6-18 * * *"
	} else {
		schedule = cr
	}

	_, err := s.Cron(schedule).Do(getTweets)
	_, _ = s.Every(1).Day().Do(webhook.GetSchedules)

	if err != nil {
		fmt.Println("Cron err", err)
		return
	}

	s.StartAsync()
}

func getTweets() {
	twts, err := tweets.TweetRepo.GetPending()

	if err != nil {
		fmt.Println("Scheduler:", err)
		return
	}

	if ShouldPost(twts[0]) {
		var isDuplicate bool

		tw := twts[0]
		fmt.Println("Posting tweet:", tw.Message)
		postErr := twitter.CreateTweet(tw.Message)

		if postErr != nil {
			fmt.Println("Posting error: ", postErr)
			isDuplicate = strings.Contains(postErr.Error(), "187")
		}

		if postErr == nil || isDuplicate {
			if isDuplicate {
				fmt.Println("Marking duplicate as posted")
			} else {
				fmt.Println("Tweeted", isDuplicate, tw.Message)
			}

			tw.Status = tweets.Posted
			tw.Modified = time.Now().Local()
			_, upErr := tweets.TweetRepo.Update(&tw)

			if upErr != nil {
				fmt.Println("error updating tweeted entry", upErr.Error(), upErr.Message())
			}
		}
	}
}

func ShouldPost(tweet tweets.Tweet) bool {
	now := time.Now().Local()

	return now.After(tweet.PostTime.Local())
}
