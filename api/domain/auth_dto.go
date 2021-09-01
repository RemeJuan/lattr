package domain

import (
	"strings"
	"time"

	"github.com/RemeJuan/lattr/utils/error_utils"
)

var scopes = []string{"token:create", "token:update", "token:read", "token:delete", "tweet:create", "tweet:update", "tweet:read", "tweet:delete"}

type Token struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Scopes    []string  `json:"scopes"`
	ExpiresAt time.Time `json:"expiry"`
	CreatedAt time.Time `json:"createdAt"`
	Modified  time.Time `json:"modified"`
	Validity  int       `json:"-"`
}

func (t *Token) Validate() error_utils.MessageErr {
	t.Name = strings.TrimSpace(t.Name)

	if t.Name == "" {
		return error_utils.UnprocessableEntityError("Name cannot be empty")
	}

	if len(t.Scopes) == 0 {
		return error_utils.UnprocessableEntityError("Scopes must be defined")
	}

	for _, val := range t.Scopes {
		if !contains(scopes, val) {
			return error_utils.UnprocessableEntityError("One or more scopes are invalid")
		}
	}

	return nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
