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
		req, err := http.NewRequest(http.MethodPost, "/tweets", bytes.NewBufferString(inputJson))
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
		path := fmt.Sprintf("/tweets/%v", recordId)
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
}
