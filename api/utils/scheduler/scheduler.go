package scheduler

import (
	"fmt"
	"os"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/twitter"
	"github.com/go-co-op/gocron"
)

func Scheduler() {
	var schedule string
	s := gocron.NewScheduler(time.UTC)
	cr := os.Getenv("CRON_SCHEDULE")

	if len(cr) == 0 {
		schedule = "*/5 6-18 * * *"
	}

	do, err := s.Cron(schedule).Do(getTweets)

	if err != nil {
		fmt.Println("Cron err", err)
		return
	}

	fmt.Println(do)

	s.StartAsync()
}

func getTweets() {
	tweets, err := domain.TweetRepo.GetPending()

	if err != nil {
		fmt.Println("Scheduler:", err)
		return
	}

	if ShouldPost(tweets[0]) {
		tw := tweets[0]
		twitter.CreateTweet(tw.Message)

		tw.Status = domain.Posted
		_, _ = domain.TweetRepo.Update(&tw)
	}
}

func ShouldPost(tweet domain.Tweet) bool {
	now := time.Now().UTC()

	return now.After(tweet.PostTime.UTC())
}
