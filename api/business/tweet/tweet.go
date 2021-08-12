package tweet

import (
	"time"

	"github.com/jinzhu/gorm"
)

type tweetStatus string

const (
	Pending = tweetStatus("Pending")
	Posted  = tweetStatus("Posted")
	Failed  = tweetStatus("Failed")
)

type Tweet struct {
	gorm.Model
	Message  string      `json:"Message"`
	UserId   string      `json:"UserId"`
	PostTime string      `json:"PostTime"`
	Status   tweetStatus `json:"Status"`
}

func ScheduleTweet(db *gorm.DB, tweet map[string]string) error {
	t, err := BuildTweet(tweet)

	if err != nil {
		return err
	}

	return db.Create(&t).Error
}

func BuildTweet(data map[string]string) (Tweet, error) {
	message := data["message"]
	postTime := data["time"]
	layout := "2006-01-02 15:04:05 -0700"

	t, err := time.Parse(layout, postTime)

	if err != nil {
		return Tweet{}, err
	}

	return Tweet{
		Message:  message,
		UserId:   "001",
		PostTime: t.UTC().String(),
		Status:   Pending,
	}, nil
}
