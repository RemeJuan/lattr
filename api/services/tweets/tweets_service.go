package tweets

import (
	"time"

	"github.com/RemeJuan/lattr/domain/tweets"
	"github.com/RemeJuan/lattr/utils/error_utils"
)

var (
	TweetService tweetServiceInterface = &tweetService{}
)

type tweetService struct{}

type tweetServiceInterface interface {
	Create(*tweets.Tweet) (*tweets.Tweet, error_utils.MessageErr)
	Get(int64) (*tweets.Tweet, error_utils.MessageErr)
	GetAll(string) ([]tweets.Tweet, error_utils.MessageErr)
	Update(*tweets.Tweet) (*tweets.Tweet, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	GetPending() ([]tweets.Tweet, error_utils.MessageErr)
	GetLast() (*tweets.Tweet, error_utils.MessageErr)
}

func (ts tweetService) Create(tweet *tweets.Tweet) (*tweets.Tweet, error_utils.MessageErr) {
	if err := tweet.Validate(); err != nil {
		return nil, err
	}

	tweet.CreatedAt = time.Now().Local()
	tweet.Modified = time.Now().Local()
	tweet.PostTime = tweet.PostTime.Local()

	tw, err := tweets.TweetRepo.Create(tweet)
	if err != nil {
		return nil, err
	}

	return tw, nil
}

func (ts tweetService) Get(id int64) (*tweets.Tweet, error_utils.MessageErr) {
	message, err := tweets.TweetRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ts tweetService) GetAll(userId string) ([]tweets.Tweet, error_utils.MessageErr) {
	messages, err := tweets.TweetRepo.GetAll(userId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ts tweetService) Update(tweet *tweets.Tweet) (*tweets.Tweet, error_utils.MessageErr) {
	if err := tweet.Validate(); err != nil {
		return nil, err
	}
	current, err := tweets.TweetRepo.Get(tweet.Id)
	if err != nil {
		return nil, err
	}
	current.Message = tweet.Message
	current.PostTime = tweet.PostTime
	current.Status = tweet.Status
	current.Modified = time.Now().Local()

	updateMsg, err := tweets.TweetRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updateMsg, nil
}

func (ts tweetService) Delete(id int64) error_utils.MessageErr {
	msg, err := tweets.TweetRepo.Get(id)
	if err != nil {
		return err
	}
	deleteErr := tweets.TweetRepo.Delete(msg.Id)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (ts tweetService) GetPending() ([]tweets.Tweet, error_utils.MessageErr) {
	messages, err := tweets.TweetRepo.GetPending()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ts tweetService) GetLast() (*tweets.Tweet, error_utils.MessageErr) {
	messages, err := tweets.TweetRepo.GetLast()
	if err != nil {
		return nil, err
	}
	return messages, nil
}
