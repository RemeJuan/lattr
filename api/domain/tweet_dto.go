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
	PostTime  time.Time   `json:"postTime"`
	Status    tweetStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	Modified  time.Time   `json:"modified"`
}

func (t *Tweet) Validate() error_utils.MessageErr {
	t.Message = strings.TrimSpace(t.Message)

	if t.Message == "" {
		return error_utils.UnprocessableEntityError("Body cannot be empty")
	}

	return nil
}
