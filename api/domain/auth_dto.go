package domain

import (
	"time"
)

type Token struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	Modified  time.Time `json:"modified"`
}
