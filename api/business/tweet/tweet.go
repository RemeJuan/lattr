package tweet

import (
	"fmt"
	"time"

	"github.com/RemeJuan/lattr/infrastructure/postgress-db"
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

func ScheduleTweet(tweet Tweet) {
	db := postgress_db.Connect()

	defer db.Close()

	db.AutoMigrate(&Tweet{})

	db.Create(&tweet)

	fmt.Println(tweet)
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
