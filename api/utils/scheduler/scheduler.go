package scheduler

import (
	"fmt"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/twitter"
	"github.com/go-co-op/gocron"
)

func Scheduler() {
	s := gocron.NewScheduler(time.UTC)

	do, err := s.Every(5).Minutes().Do(getTweets)

	if err != nil {
		fmt.Println("Cron err", err)
		return
	}

	fmt.Println(do)

	//s.StartAsync()
}

func getTweets() {
	tweets, err := domain.TweetRepo.GetPending()

	if err != nil {
		fmt.Println("Scheduler:", err)
		return
	}

	if ShouldPost(tweets[0]) {
		twitter.CreateTweet(tweets[0].Message)
	}
}

func ShouldPost(tweet domain.Tweet) bool {
	now := time.Now().UTC()

	return now.After(tweet.PostTime.UTC())
}
