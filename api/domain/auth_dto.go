package domain

import (
	"strings"
	"time"

	"github.com/RemeJuan/lattr/utils/error_utils"
)

type Token struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	Modified  time.Time `json:"modified"`
}

func (t *Token) Validate() error_utils.MessageErr {
	t.Name = strings.TrimSpace(t.Name)

	if t.Name == "" {
		return error_utils.UnprocessableEntityError("Name cannot be empty")
	}

	return nil
}
