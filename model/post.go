package model

import (
	"time"
)

type Post struct {
	Id           string    `json:"id"`
	Message      string    `json:"message"`
	CreationDate time.Time `json:"creationDate"`
}
