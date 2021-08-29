package controllers

import (
	"net/http"
	"strconv"
	"time"

	tweets2 "github.com/RemeJuan/lattr/domain/tweets"
	"github.com/RemeJuan/lattr/services/tweets"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/RemeJuan/lattr/utils/webhook"
	"github.com/gin-gonic/gin"
)

func CreateTweet(c *gin.Context) {
	var tweet tweets2.Tweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	tweet.Status = tweets2.Pending
	msg, err := tweets.TweetService.Create(&tweet)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func GetTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	message, getErr := tweets.TweetService.Get(twId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, message)
}

func GetTweets(c *gin.Context) {
	userId := GetParam(c, "userId")

	tweets, getErr := tweets.TweetService.GetAll(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, tweets)
}

func UpdateTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}
	var tweet tweets2.Tweet
	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	tweet.Id = twId
	msg, err := tweets.TweetService.Update(&tweet)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, msg)
}

func DeleteTweet(c *gin.Context) {
	paramId := GetParam(c, "id")
	twId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	if err := tweets.TweetService.Delete(twId); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func WebHook(c *gin.Context) {
	var tweet tweets2.Tweet
	var tweetTime time.Time

	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	last, lErr := tweets.TweetService.GetLast()
	if last == nil {
		tweetTime = webhook.DetermineScheduleType(time.Now().Local())
	} else {
		tweetTime = webhook.DetermineScheduleType(last.PostTime)
	}

	if lErr != nil && lErr.Error() != "not_found" {
		c.JSON(lErr.Status(), lErr)
		return
	}

	tweet.PostTime = tweetTime
	tweet.Status = tweets2.Scheduled
	msg, err := tweets.TweetService.Create(&tweet)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func GetParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
