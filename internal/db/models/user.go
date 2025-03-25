package models

import (
	"time"
)

type User struct {
    ChatID      int64
    CreatedAt   time.Time
    Name        string
}
