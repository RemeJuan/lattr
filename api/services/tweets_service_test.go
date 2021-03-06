package services

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/stretchr/testify/assert"
)

var (
	createTweetDomain      func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	getTweetDomain         func(messageId int64) (*domain.Tweet, error_utils.MessageErr)
	updateTweetDomain      func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	deleteTweetDomain      func(messageId int64) error_utils.MessageErr
	getAllTweetsDomain     func(userId string) ([]domain.Tweet, error_utils.MessageErr)
	getPendingTweetsDomain func() ([]domain.Tweet, error_utils.MessageErr)
	getLastTweetsDomain    func() (*domain.Tweet, error_utils.MessageErr)
)

type tweetDbMock struct {
	domain.TweetRepoInterface
}

func (m *tweetDbMock) Create(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return createTweetDomain(msg)
}
func (m *tweetDbMock) Get(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
	return getTweetDomain(messageId)
}
func (m *tweetDbMock) GetAll(userId string) ([]domain.Tweet, error_utils.MessageErr) {
	return getAllTweetsDomain(userId)
}
func (m *tweetDbMock) Update(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return updateTweetDomain(msg)
}
func (m *tweetDbMock) Delete(messageId int64) error_utils.MessageErr {
	return deleteTweetDomain(messageId)
}
func (m *tweetDbMock) GetPending() ([]domain.Tweet, error_utils.MessageErr) {
	return getPendingTweetsDomain()
}
func (m *tweetDbMock) GetLast() (*domain.Tweet, error_utils.MessageErr) {
	return getLastTweetsDomain()
}
func (m *tweetDbMock) Initialize() *sql.DB {
	return nil
}

const layout = "2021-07-12 10:55:50 +0000"

func TestTweetService_Create(t *testing.T) {
	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = "the message"

		createTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		request := &domain.Tweet{
			Message:   message,
			PostTime:  postTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Create(request)

		assert.Nil(t, err)
		assert.NotNil(t, msg)
		assert.EqualValues(t, recordId, msg.Id)
		assert.EqualValues(t, message, msg.Message)
		assert.EqualValues(t, postTime, msg.PostTime)
		assert.EqualValues(t, tm, msg.CreatedAt)
	})

	t.Run("Validation failed", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = ""

		createTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		request := &domain.Tweet{
			Message:   message,
			PostTime:  postTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Create(request)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, err.Status())
		assert.EqualValues(t, "invalid_request", err.Error())
		assert.EqualValues(t, "Body cannot be empty", err.Message())
	})

	t.Run("Create failed", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = "message"

		createTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.InternalServerError("Unknown error occurred")
		}

		request := &domain.Tweet{
			Message:   message,
			PostTime:  postTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Create(request)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
		assert.EqualValues(t, "Unknown error occurred", err.Message())
	})
}

func TestTweetService_Get(t *testing.T) {
	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")
	var message = ""

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		message = "the message"

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		tweet, err := TweetService.Get(1)

		assert.Nil(t, err)
		assert.NotNil(t, tweet)
		assert.EqualValues(t, recordId, tweet.Id)
		assert.EqualValues(t, message, tweet.Message)
		assert.EqualValues(t, postTime, tweet.PostTime)
		assert.EqualValues(t, tm, tweet.CreatedAt)
	})

	t.Run("Not Found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("the id is not found")
		}

		msg, err := TweetService.Get(recordId)
		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "the id is not found", err.Message())
		assert.EqualValues(t, "not_found", err.Error())
	})
}

func TestTweetService_GetAll(t *testing.T) {
	const userId = "001"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")
	var message = ""

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		message = "the message"

		getAllTweetsDomain = func(userId string) ([]domain.Tweet, error_utils.MessageErr) {
			return []domain.Tweet{
				{
					Id:        01,
					Message:   message,
					PostTime:  postTime,
					CreatedAt: tm,
					UserId:    userId,
				},
				{
					Id:        02,
					Message:   message,
					PostTime:  postTime,
					CreatedAt: tm,
					UserId:    userId,
				},
			}, nil
		}

		tweets, err := TweetService.GetAll(userId)

		assert.Nil(t, err)
		assert.NotNil(t, tweets)

		// First Result
		assert.EqualValues(t, 01, tweets[0].Id)
		assert.EqualValues(t, userId, tweets[0].UserId)
		assert.EqualValues(t, message, tweets[0].Message)
		assert.EqualValues(t, postTime, tweets[0].PostTime)

		// Second Result
		assert.EqualValues(t, 02, tweets[1].Id)
		assert.EqualValues(t, userId, tweets[1].UserId)
		assert.EqualValues(t, message, tweets[1].Message)
		assert.EqualValues(t, postTime, tweets[1].PostTime)
	})

	t.Run("Not Found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getAllTweetsDomain = func(userId string) ([]domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("error getting messages")
		}

		msg, err := TweetService.GetAll(userId)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "error getting messages", err.Message())
		assert.EqualValues(t, "not_found", err.Error())
	})
}

func TestTweetService_Update(t *testing.T) {
	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")
	updatedTime, _ := time.Parse(layout, "2021-07-13 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = "the message"

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		updateTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  updatedTime,
				CreatedAt: tm,
			}, nil
		}

		request := &domain.Tweet{
			Message:   message,
			PostTime:  updatedTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Update(request)

		assert.Nil(t, err)
		assert.NotNil(t, msg)
		assert.EqualValues(t, recordId, msg.Id)
		assert.EqualValues(t, message, msg.Message)
		assert.EqualValues(t, updatedTime, msg.PostTime)
		assert.EqualValues(t, tm, msg.CreatedAt)
	})

	t.Run("Validation failed", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = "the message"
		const updatedMessage = ""

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		updateTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   updatedMessage,
				PostTime:  updatedTime,
				CreatedAt: tm,
			}, nil
		}

		request := &domain.Tweet{
			Message:   updatedMessage,
			PostTime:  updatedTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Update(request)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, err.Status())
		assert.EqualValues(t, "invalid_request", err.Error())
		assert.EqualValues(t, "Body cannot be empty", err.Message())
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const updatedMessage = "message"
		const errMessage = "No matching record found"

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError(errMessage)
		}

		request := &domain.Tweet{
			Message:   updatedMessage,
			PostTime:  updatedTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Update(request)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "not_found", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})

	t.Run("Update failed", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		const message = "the message"
		const errMessage = "Unable to reset token"

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		updateTweetDomain = func(msg *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.InternalServerError(errMessage)
		}

		request := &domain.Tweet{
			Message:   message,
			PostTime:  updatedTime,
			CreatedAt: tm,
		}
		msg, err := TweetService.Update(request)

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})
}

func TestTweetService_Delete(t *testing.T) {
	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")
	var message = ""

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		message = "the message"

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:        recordId,
				Message:   message,
				PostTime:  postTime,
				CreatedAt: tm,
			}, nil
		}

		deleteTweetDomain = func(messageId int64) error_utils.MessageErr {
			return nil
		}

		err := TweetService.Delete(recordId)

		assert.Nil(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}
		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("Something went wrong getting message")
		}
		err := TweetService.Delete(1)
		assert.NotNil(t, err)
		assert.EqualValues(t, "Something went wrong getting message", err.Message())
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "not_found", err.Error())
	})

	t.Run("Unable to delete", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getTweetDomain = func(messageId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       1,
				Message:  message,
				PostTime: postTime,
			}, nil
		}

		deleteTweetDomain = func(messageId int64) error_utils.MessageErr {
			return error_utils.InternalServerError("error deleting message")
		}

		err := TweetService.Delete(recordId)

		assert.NotNil(t, err)
		assert.EqualValues(t, "error deleting message", err.Message())
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
	})
}

func TestTweetService_GetPending(t *testing.T) {
	const userId = "001"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")
	var message = ""

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		message = "the message"

		getPendingTweetsDomain = func() ([]domain.Tweet, error_utils.MessageErr) {
			return []domain.Tweet{
				{
					Id:        01,
					Message:   message,
					PostTime:  postTime,
					CreatedAt: tm,
					UserId:    userId,
				},
				{
					Id:        02,
					Message:   message,
					PostTime:  postTime,
					CreatedAt: tm,
					UserId:    userId,
				},
			}, nil
		}

		tweets, err := TweetService.GetPending()

		assert.Nil(t, err)
		assert.NotNil(t, tweets)

		// First Result
		assert.EqualValues(t, 01, tweets[0].Id)
		assert.EqualValues(t, userId, tweets[0].UserId)
		assert.EqualValues(t, message, tweets[0].Message)
		assert.EqualValues(t, postTime, tweets[0].PostTime)

		// Second Result
		assert.EqualValues(t, 02, tweets[1].Id)
		assert.EqualValues(t, userId, tweets[1].UserId)
		assert.EqualValues(t, message, tweets[1].Message)
		assert.EqualValues(t, postTime, tweets[1].PostTime)
	})

	t.Run("Not Found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getPendingTweetsDomain = func() ([]domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("error getting messages")
		}

		msg, err := TweetService.GetPending()

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "error getting messages", err.Message())
		assert.EqualValues(t, "not_found", err.Error())
	})
}

func TestTweetService_GetLast(t *testing.T) {
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getLastTweetsDomain = func() (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				PostTime: postTime,
			}, nil
		}

		tweet, err := TweetService.GetLast()

		assert.Nil(t, err)
		assert.NotNil(t, tweet)

		// First Result
		assert.EqualValues(t, postTime, tweet.PostTime)
	})

	t.Run("Not Found", func(t *testing.T) {
		domain.TweetRepo = &tweetDbMock{}

		getLastTweetsDomain = func() (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("error getting messages")
		}

		msg, err := TweetService.GetLast()

		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "error getting messages", err.Message())
		assert.EqualValues(t, "not_found", err.Error())
	})
}
