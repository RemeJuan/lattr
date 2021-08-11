package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RemeJuan/lattr/business/tweet"
	"github.com/RemeJuan/lattr/infrastructure/twitter-client"
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

	InternalServerError(c, err)

	if len(webhookData["data"]) > 0 {
		twitter_client.CreateTweet(webhookData["data"])
	} else {
		BadRequestError(c, errors.New("invalid payload"))
	}
}

func (con *DBController) Tweet(c *gin.Context) {
	payload := make(map[string]string)
	err := json.NewDecoder(c.Request.Body).Decode(&payload)

	PingDB(con.db)

	InternalServerError(c, err)

	tweet.ScheduleTweet(con.db, payload)
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

func InternalServerError(c *gin.Context, err error) {
	http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	return
}

func BadRequestError(c *gin.Context, err error) {
	http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	return
}

func Decoder(c *gin.Context, data *map[string]string) error {
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	return err
}
