package models

import (
	"p2p-stats/internal/domain/entities"
	"time"
)

type User struct {
    ChatID      int64
    CreatedAt   time.Time
    Name        string
}

func ToDBUser(user *entities.ValidatedUser) *User {
    return &User{
        ChatID: user.ChatID,
        CreatedAt: user.CreatedAt,
        Name: user.Name,
    }
}

func FromDBUser(user *User) *entities.User {
    return &entities.User{
        ChatID: user.ChatID,
        CreatedAt: user.CreatedAt,
        Name: user.Name,
    }
}
