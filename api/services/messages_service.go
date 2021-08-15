package services

import (
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
)

var (
	TweetService tweetServiceInterface = &tweetService{}
)

type tweetService struct{}

type tweetServiceInterface interface {
	Create(*domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	Get(int64) (*domain.Tweet, error_utils.MessageErr)
	GetAll(string) ([]domain.Tweet, error_utils.MessageErr)
	Update(*domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
}

func (ts tweetService) Create(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	if err := tweet.Validate(); err != nil {
		return nil, err
	}

	tweet.CreatedAt = time.Now()
	tweet, err := domain.TweetRepo.Create(tweet)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (ts tweetService) Get(id int64) (*domain.Tweet, error_utils.MessageErr) {
	message, err := domain.TweetRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ts tweetService) GetAll(userId string) ([]domain.Tweet, error_utils.MessageErr) {
	messages, err := domain.TweetRepo.GetAll(userId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ts tweetService) Update(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	if err := tweet.Validate(); err != nil {
		return nil, err
	}
	current, err := domain.TweetRepo.Get(tweet.Id)
	if err != nil {
		return nil, err
	}
	current.Message = tweet.Message
	current.PostTime = tweet.PostTime
	current.Status = tweet.Status
	current.Modified = time.Now()

	updateMsg, err := domain.TweetRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updateMsg, nil
}

func (ts tweetService) Delete(id int64) error_utils.MessageErr {
	msg, err := domain.TweetRepo.Get(id)
	if err != nil {
		return err
	}
	deleteErr := domain.TweetRepo.Delete(msg.Id)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}
