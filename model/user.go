package model

import (
	"time"
)

type User struct {
	ID           int
	Email        string
	Password     string
	Role         string
	Active       bool
	Created      time.Time
	Subscription time.Time
}

func NewUser() User {
	user := &User{}
	user.Role = "normal"
	user.Active = true
	return *user
}
