package scheduler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/twitter"
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

	if err != nil {
		fmt.Println("Cron err", err)
		return
	}

	s.StartAsync()
}

func getTweets() {
	tweets, err := domain.TweetRepo.GetPending()

	if err != nil {
		fmt.Println("Scheduler:", err)
		return
	}

	if ShouldPost(tweets[0]) {
		var isDuplicate bool

		tw := tweets[0]
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

			tw.Status = domain.Posted
			tw.Modified = time.Now().Local()
			_, upErr := domain.TweetRepo.Update(&tw)

			if upErr != nil {
				fmt.Println("error updating tweeted entry", upErr.Error(), upErr.Message())
			}
		}
	}
}

func ShouldPost(tweet domain.Tweet) bool {
	now := time.Now().Local()

	return now.After(tweet.PostTime.Local())
}
