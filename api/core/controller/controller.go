package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RemeJuan/lattr/business/tweet"
	twitter_client "github.com/RemeJuan/lattr/infrastructure/twitter-client"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DBController struct {
	db *gorm.DB
}

func Init(db *gorm.DB) *DBController {
	return &DBController{
		db: db,
	}
}

func (con *DBController) WebHook(c *gin.Context) {
	webhookData := make(map[string]string)
	err := Decoder(c, &webhookData)

	UnprocessableEntity(c, err)

	if len(webhookData["data"]) > 0 {
		twitter_client.CreateTweet(webhookData["data"])
		c.JSON(http.StatusCreated, "Tweet Queued")
	} else {
		BadRequestError(c, errors.New("invalid payload"))
	}
}

func (con *DBController) Tweet(c *gin.Context) {
	payload := make(map[string]string)
	err := Decoder(c, &payload)

	PingDB(con.db)

	InternalServerError(c, err)

	tErr := tweet.ScheduleTweet(con.db, payload)

	ErrorCheck(tErr)

	c.JSON(200, "Tweet Queued")
}

func PingDB(con *gorm.DB) {
	err := con.DB().Ping()
	ErrorCheck(err)
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func UnprocessableEntity(c *gin.Context, err error) {
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}
}

func InternalServerError(c *gin.Context, err error) {
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func BadRequestError(c *gin.Context, err error) {
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func Decoder(c *gin.Context, data *map[string]string) error {
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	return err
}
