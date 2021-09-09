package controllers

import (
	"net/http"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/RemeJuan/lattr/utils/webhook"
	"github.com/gin-gonic/gin"
)

func WebHook(c *gin.Context) {
	var tweet domain.Tweet
	var tweetTime time.Time

	if err := c.ShouldBindJSON(&tweet); err != nil {
		theErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	last, lErr := services.TweetService.GetLast()
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
	tweet.Status = domain.Scheduled
	msg, err := services.TweetService.Create(&tweet)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}
