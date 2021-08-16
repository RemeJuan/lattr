package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	createTweetService func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	getTweetService    func(msgId int64) (*domain.Tweet, error_utils.MessageErr)
	updateTweetService func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	deleteTweetService func(msgId int64) error_utils.MessageErr
	getAllTweetService func(userId string) ([]domain.Tweet, error_utils.MessageErr)
)

type serviceMock struct {
}

func (sm *serviceMock) Create(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return createTweetService(tweet)
}

func (sm *serviceMock) Get(id int64) (*domain.Tweet, error_utils.MessageErr) {
	return getTweetService(id)
}

func (sm *serviceMock) GetAll(userId string) ([]domain.Tweet, error_utils.MessageErr) {
	return getAllTweetService(userId)
}

func (sm *serviceMock) Update(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return updateTweetService(tweet)
}

func (sm *serviceMock) Delete(id int64) error_utils.MessageErr {
	return deleteTweetService(id)
}

const basePath = "/tweets"

func TestCreateTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1
	const postTime = "2021-07-12 10:55:50 +0000"

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		const msg = "the message"

		createTweetService = func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  msg,
				PostTime: postTime,
			}, nil
		}
		jsonBody := `{"title": "the title", "body": "the body"}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, basePath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(basePath, CreateTweet)
		r.ServeHTTP(rr, req)

		var tweet domain.Tweet
		err = json.Unmarshal(rr.Body.Bytes(), &tweet)
		assert.Nil(t, err)
		assert.NotNil(t, tweet)
		assert.EqualValues(t, http.StatusCreated, rr.Code)
		assert.EqualValues(t, recordId, tweet.Id)
		assert.EqualValues(t, msg, tweet.Message)
		assert.EqualValues(t, postTime, tweet.PostTime)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		inputJson := ""
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, basePath, bytes.NewBufferString(inputJson))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST("/tweets", CreateTweet)
		r.ServeHTTP(rr, req)

		apiErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
		assert.EqualValues(t, "invalid json body", apiErr.Message())
		assert.EqualValues(t, "invalid_request", apiErr.Error())
	})

	t.Run("Create Error", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		createTweetService = func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.UnprocessableEntityError("Body cannot be empty")
		}
		jsonBody := `{}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, basePath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(basePath, CreateTweet)
		r.ServeHTTP(rr, req)

		msgErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "Body cannot be empty", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})
}

func TestGetTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1
	const postTime = "2021-07-12 10:55:50 +0000"

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		const message = "the message"

		getTweetService = func(msgId int64) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  message,
				PostTime: postTime,
				Status:   domain.Pending,
			}, nil
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", basePath, recordId)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		r.GET("/tweets/:id", GetTweet)
		r.ServeHTTP(rr, req)

		var tweet domain.Tweet
		err := json.Unmarshal(rr.Body.Bytes(), &tweet)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, recordId, tweet.Id)
		assert.EqualValues(t, message, tweet.Message)
		assert.EqualValues(t, postTime, tweet.PostTime)
		assert.EqualValues(t, domain.Pending, tweet.Status)
	})

	t.Run("Cannot parse ID", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		const invalidID = "red"

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", basePath, invalidID)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		r.GET("/tweets/:id", GetTweet)
		r.ServeHTTP(rr, req)

		msgErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "unable to parse ID", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})

	t.Run("Error", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		getTweetService = func(msgId int64) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("unable to find item")
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", basePath, recordId)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		r.GET("/tweets/:id", GetTweet)
		r.ServeHTTP(rr, req)

		msgErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusNotFound, msgErr.Status())
		assert.EqualValues(t, "unable to find item", msgErr.Message())
		assert.EqualValues(t, "not_found", msgErr.Error())
	})
}

func TestGetTweets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1
	const postTime = "2021-07-12 10:55:50 +0000"

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		const message = "the message"
		const userId = "test"

		getAllTweetService = func(userId string) ([]domain.Tweet, error_utils.MessageErr) {
			return []domain.Tweet{
				{
					Id:       recordId,
					Message:  message,
					PostTime: postTime,
					Status:   domain.Pending,
				},
				{
					Id:       recordId,
					Message:  message,
					PostTime: postTime,
					Status:   domain.Pending,
				},
			}, nil
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/all/%v", basePath, userId)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		r.GET("/tweets/all/:userId", GetTweets)
		r.ServeHTTP(rr, req)

		var tweets []domain.Tweet
		err := json.Unmarshal(rr.Body.Bytes(), &tweets)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, recordId, tweets[0].Id)
		assert.EqualValues(t, message, tweets[0].Message)
		assert.EqualValues(t, postTime, tweets[0].PostTime)
		assert.EqualValues(t, domain.Pending, tweets[0].Status)

		assert.EqualValues(t, recordId, tweets[1].Id)
		assert.EqualValues(t, message, tweets[1].Message)
		assert.EqualValues(t, postTime, tweets[1].PostTime)
		assert.EqualValues(t, domain.Pending, tweets[1].Status)
	})

	t.Run("Error", func(t *testing.T) {
		services.TweetService = &serviceMock{}

		const userId = "test"

		getAllTweetService = func(userId string) ([]domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("No results")
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/all/%v", basePath, userId)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		r.GET("/tweets/all/:userId", GetTweets)
		r.ServeHTTP(rr, req)

		apiErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.Nil(t, err)
		assert.EqualValues(t, apiErr.Status(), rr.Code)
		assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
		assert.EqualValues(t, "not_found", apiErr.Error())
		assert.EqualValues(t, "No results", apiErr.Message())
	})
}
