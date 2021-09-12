package domain

import (
	"strings"
	"time"

	"github.com/RemeJuan/lattr/utils/error_utils"
)

type tweetStatus string

const (
	Pending   = tweetStatus("Pending")
	Posted    = tweetStatus("Posted")
	Scheduled = tweetStatus("Scheduled")
)

type Tweet struct {
	Id        int64       `json:"id" example:"1"`
	Message   string      `json:"message" example:"TIL: Life is awesome"`
	UserId    string      `json:"userId" example:"IFTTT"`
	Status    tweetStatus `json:"status" example:"Pending"`
	PostTime  time.Time   `json:"postTime" example:"2022-09-09T10:29:07.559636Z"`
	CreatedAt time.Time   `json:"createdAt" example:"2022-09-09T10:29:07.559636Z"`
	Modified  time.Time   `json:"modified" example:"2022-09-09T10:29:07.559636Z"`
}

func (t *Tweet) Validate() error_utils.MessageErr {
	t.Message = strings.TrimSpace(t.Message)

	if t.Message == "" {
		return error_utils.UnprocessableEntityError("Body cannot be empty")
	}

	return nil
}
