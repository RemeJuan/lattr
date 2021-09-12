package controllers

import (
	"net/http"
	"strconv"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
)

// CreateTweet godoc
// @Summary Create a new tweet
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param tweet body domain.Tweet true "Create tweet"
// @Success 200 {object} domain.Tweet
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[tweet:create]
// @Router /tweets/create [post]
func CreateTweet(c *gin.Context) {
	var tweet domain.Tweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	tweet.Status = domain.Pending
	msg, err := services.TweetService.Create(&tweet)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}

// GetTweet godoc
// @Summary Fetch Tweet by ID
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param id path int true "Tweet ID"
// @Success 200 {object} domain.Tweet
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Router /tweets/{id} [get]
func GetTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	message, getErr := services.TweetService.Get(twId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, message)
}

// GetTweets godoc
// @Summary List all Tweets by UserId
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param userId path int true "User ID"
// @Success 200 {array} domain.Tweet
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Router /tweets/all/{userId} [get]
func GetTweets(c *gin.Context) {
	userId := GetParam(c, "userId")

	tweets, getErr := services.TweetService.GetAll(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, tweets)
}

// UpdateTweet godoc
// @Summary Updated a single tweet
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param id path int true "Tweet ID"
// @Param tweet body domain.Tweet true "Update tweet"
// @Success 200 {object} domain.Tweet
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Router /tweets/{id} [put]
func UpdateTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}
	var tweet domain.Tweet
	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	tweet.Id = twId
	msg, err := services.TweetService.Update(&tweet)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, msg)
}

// DeleteTweet godoc
// @Summary Deletes a single tweet
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param id path int true "Tweet ID"
// @Success 200 {object} object "{message: "success"}"
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Router /tweets/{id} [delete]
func DeleteTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	if err := services.TweetService.Delete(twId); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
