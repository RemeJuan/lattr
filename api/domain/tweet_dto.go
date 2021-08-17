package domain

import (
	"strings"
	"time"

	"github.com/RemeJuan/lattr/utils/error_utils"
)

type tweetStatus string

const (
	Pending = tweetStatus("Pending")
	Posted  = tweetStatus("Posted")
	Failed  = tweetStatus("Failed")
)

type Tweet struct {
	Id        int64       `json:"id"`
	Message   string      `json:"message"`
	UserId    string      `json:"userId"`
	PostTime  string      `json:"postTime"`
	Status    tweetStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	Modified  time.Time   `json:"modified"`
}

func (t *Tweet) Validate() error_utils.MessageErr {
	t.Message = strings.TrimSpace(t.Message)
	t.PostTime = strings.TrimSpace(t.PostTime)
	layout := "2006-01-02 15:04:05 -0700"
	_, err := time.Parse(layout, t.PostTime)

	if t.Message == "" {
		return error_utils.UnprocessableEntityError("Body cannot be empty")
	}

	if err != nil {
		return error_utils.UnprocessableEntityError(`Invalid date/time format, please use "CCYY-MM-DD HH:mm:ss -zzzz`)
	}

	return nil
}
