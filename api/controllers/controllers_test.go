package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/RemeJuan/lattr/utils/webhook"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const tweetPath = "/tweets"
const tokenPath = "/token"
const layout = "2021-07-12 10:55:50 +0000"

func TestCreateTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

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
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(tweetPath, CreateTweet)
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
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(inputJson))
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
		services.TweetService = &tweetServiceMock{}

		createTweetService = func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.UnprocessableEntityError("Body cannot be empty")
		}
		jsonBody := `{}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(tweetPath, CreateTweet)
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
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

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
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
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
		services.TweetService = &tweetServiceMock{}

		const invalidID = "red"

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, invalidID)
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
		services.TweetService = &tweetServiceMock{}

		getTweetService = func(msgId int64) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("unable to find item")
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
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
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

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
		path := fmt.Sprintf("%s/all/%v", tweetPath, userId)
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
		services.TweetService = &tweetServiceMock{}

		const userId = "test"

		getAllTweetService = func(userId string) ([]domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("No results")
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/all/%v", tweetPath, userId)
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

func TestUpdateTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		const message = "different message"

		updateTweetService = func(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  message,
				PostTime: postTime,
				Status:   domain.Pending,
			}, nil
		}
		inputJson := `{"message": "different message"}`
		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
		req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBufferString(inputJson))
		rr := httptest.NewRecorder()
		r.PUT("/tweets/:id", UpdateTweet)
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
		services.TweetService = &tweetServiceMock{}

		const invalidID = "red"

		inputJson := `{"message": "different message"}`
		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, invalidID)
		req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBufferString(inputJson))
		rr := httptest.NewRecorder()
		r.PUT("/tweets/:id", UpdateTweet)
		r.ServeHTTP(rr, req)

		msgErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "unable to parse ID", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		inputJson := ""
		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
		req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBufferString(inputJson))
		rr := httptest.NewRecorder()
		r.PUT("/tweets/:id", UpdateTweet)
		r.ServeHTTP(rr, req)

		msgErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "invalid json body", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})

	t.Run("Error", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		updateTweetService = func(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError("unable to find item")
		}

		inputJson := `{"message": "different message"}`
		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
		req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBufferString(inputJson))
		rr := httptest.NewRecorder()
		r.PUT("/tweets/:id", UpdateTweet)
		r.ServeHTTP(rr, req)

		msgErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusNotFound, msgErr.Status())
		assert.EqualValues(t, "unable to find item", msgErr.Message())
		assert.EqualValues(t, "not_found", msgErr.Error())
	})
}

func TestDeleteTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const recordId int64 = 1

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		deleteTweetService = func(id int64) error_utils.MessageErr {
			return nil
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
		req, _ := http.NewRequest(http.MethodDelete, path, nil)
		rr := httptest.NewRecorder()
		r.DELETE("/tweets/:id", DeleteTweet)
		r.ServeHTTP(rr, req)

		var result map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &result)

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.Equal(t, "deleted", result["status"])
	})

	t.Run("Unable to parse ID", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, "red")
		req, _ := http.NewRequest(http.MethodDelete, path, nil)
		rr := httptest.NewRecorder()
		r.DELETE("/tweets/:id", DeleteTweet)
		r.ServeHTTP(rr, req)

		apiErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)
		assert.EqualValues(t, rr.Code, apiErr.Status())
		assert.Equal(t, "unable to parse ID", apiErr.Message())
		assert.Equal(t, "invalid_request", apiErr.Error())
	})

	t.Run("Error", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		deleteTweetService = func(id int64) error_utils.MessageErr {
			return error_utils.InternalServerError("Unable to delete entry")
		}

		r := gin.Default()
		path := fmt.Sprintf("%s/%v", tweetPath, recordId)
		req, _ := http.NewRequest(http.MethodDelete, path, nil)
		rr := httptest.NewRecorder()
		r.DELETE("/tweets/:id", DeleteTweet)
		r.ServeHTTP(rr, req)

		apiErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
		assert.EqualValues(t, rr.Code, apiErr.Status())
		assert.Equal(t, "Unable to delete entry", apiErr.Message())
		assert.Equal(t, "server_error", apiErr.Error())
	})
}

func TestWebHook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	schedules := []string{"14:30", "15:31"}
	_ = os.Setenv("SCHEDULE_TYPE", "FIXED")
	_ = os.Setenv("SCHEDULES", strings.Join(schedules, ","))

	webhook.GetSchedules()

	const recordId = 1
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		const msg = "the message"

		getLastTweet = func() (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  msg,
				PostTime: postTime,
			}, nil
		}

		createTweetService = func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  msg,
				PostTime: postTime,
			}, nil
		}

		jsonBody := `{"message": "the message"}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(tweetPath, WebHook)
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
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(inputJson))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST("/tweets", WebHook)
		r.ServeHTTP(rr, req)

		apiErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
		assert.EqualValues(t, "invalid json body", apiErr.Message())
		assert.EqualValues(t, "invalid_request", apiErr.Error())
	})

	t.Run("GetLast Error", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		getLastTweet = func() (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.UnprocessableEntityError("No tweets found")
		}

		jsonBody := `{}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(tweetPath, WebHook)
		r.ServeHTTP(rr, req)

		msgErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "No tweets found", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})

	t.Run("Create Error", func(t *testing.T) {
		services.TweetService = &tweetServiceMock{}

		const msg = "the message"

		getLastTweet = func() (*domain.Tweet, error_utils.MessageErr) {
			return &domain.Tweet{
				Id:       recordId,
				Message:  msg,
				PostTime: postTime,
			}, nil
		}

		createTweetService = func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
			return nil, error_utils.UnprocessableEntityError("Body cannot be empty")
		}
		jsonBody := `{}`
		r := gin.Default()
		req, err := http.NewRequest(http.MethodPost, tweetPath, bytes.NewBufferString(jsonBody))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.POST(tweetPath, CreateTweet)
		r.ServeHTTP(rr, req)

		msgErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

		assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
		assert.EqualValues(t, "Body cannot be empty", msgErr.Message())
		assert.EqualValues(t, "invalid_request", msgErr.Error())
	})
}

func TestAuthControllers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const mockId int64 = 1
	const mockName = "token"
	const mockToken = "mock-token"
	mockScopes := []string{"token:create"}
	mockDate := time.Date(2021, 8, 20, 14, 30, 0, 0, time.Local)

	t.Run("CreateToken", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			services.AuthService = &authServiceMock{}

			createTokenService = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
				return &domain.Token{
					Id:        mockId,
					Name:      mockName,
					Token:     mockToken,
					Scopes:    mockScopes,
					ExpiresAt: mockDate,
					CreatedAt: mockDate,
					Modified:  mockDate,
				}, nil
			}
			jsonBody := `{"name": "token", "scope": ["token:create"]}`
			r := gin.Default()
			req, err := http.NewRequest(http.MethodPost, tokenPath, bytes.NewBufferString(jsonBody))
			if err != nil {
				t.Errorf("this is the error: %v\n", err)
			}
			rr := httptest.NewRecorder()
			r.POST(tokenPath, CreateToken)
			r.ServeHTTP(rr, req)

			var token domain.Token
			err = json.Unmarshal(rr.Body.Bytes(), &token)
			assert.Nil(t, err)
			assert.NotNil(t, token)
			assert.EqualValues(t, http.StatusCreated, rr.Code)
			assert.EqualValues(t, mockId, token.Id)
			assert.EqualValues(t, mockName, token.Name)
			assert.EqualValues(t, mockScopes, token.Scopes)
		})

		t.Run("Invalid JSON", func(t *testing.T) {
			inputJson := ""
			r := gin.Default()
			req, err := http.NewRequest(http.MethodPost, tokenPath, bytes.NewBufferString(inputJson))
			if err != nil {
				t.Errorf("this is the error: %v\n", err)
			}
			rr := httptest.NewRecorder()
			r.POST(tokenPath, CreateToken)
			r.ServeHTTP(rr, req)

			apiErr, _ := error_utils.ApiErrFromBytes(rr.Body.Bytes())

			assert.Nil(t, err)
			assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
			assert.EqualValues(t, "invalid json body", apiErr.Message())
			assert.EqualValues(t, "invalid_request", apiErr.Error())
		})

		t.Run("Create Error", func(t *testing.T) {
			services.AuthService = &authServiceMock{}

			createTokenService = func(message *domain.Token) (*domain.Token, error_utils.MessageErr) {
				return nil, error_utils.UnprocessableEntityError("Body cannot be empty")
			}
			jsonBody := `{}`
			r := gin.Default()
			req, err := http.NewRequest(http.MethodPost, tokenPath, bytes.NewBufferString(jsonBody))
			if err != nil {
				t.Errorf("this is the error: %v\n", err)
			}
			rr := httptest.NewRecorder()
			r.POST(tokenPath, CreateTweet)
			r.ServeHTTP(rr, req)

			msgErr, err := error_utils.ApiErrFromBytes(rr.Body.Bytes())

			assert.EqualValues(t, http.StatusUnprocessableEntity, msgErr.Status())
			assert.EqualValues(t, "Body cannot be empty", msgErr.Message())
			assert.EqualValues(t, "invalid_request", msgErr.Error())
		})
	})

}
