package domain

import (
	"strings"
	"time"

	"github.com/RemeJuan/lattr/utils/error_utils"
)

var scopes = []string{"token:create", "token:update", "token:read", "token:delete", "tweet:create", "tweet:update", "tweet:read", "tweet:delete"}

type Token struct {
	Id        int64     `json:"id" example:"1"`
	Name      string    `json:"name" example:"IFTTT"`
	Token     string    `json:"token" example:"1d6dcc23-51c4-4540-b659-b2834efad5bc"`
	Scopes    []string  `json:"scopes" example:"[token:create]"`
	ExpiresAt time.Time `json:"expiry" example:"2022-09-09T10:29:07.559636Z"`
	CreatedAt time.Time `json:"createdAt" example:"2022-09-09T10:29:07.559636Z"`
	Modified  time.Time `json:"modified" example:"2022-09-09T10:29:07.559636Z"`
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
