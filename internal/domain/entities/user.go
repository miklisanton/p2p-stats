package entities

import (
	"errors"
	"time"
)

type User struct {
    ChatID      int64
    CreatedAt   time.Time
    Name        string
}

func NewUser(chatID int64, name string) *User {
    return &User{
        ChatID: chatID,
        CreatedAt: time.Now(),
        Name: name,
    }
}

func (u *User) validate() error {
    if u.Name == "" {
        return errors.New("invalid name")
    }
    return nil
}
