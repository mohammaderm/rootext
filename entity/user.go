package entity

import (
	"time"
)

type User struct {
	ID         uint
	Username   string
	Password   string
	Created_at time.Time
}
