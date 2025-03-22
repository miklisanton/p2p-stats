package telegram

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
    Bot         *tgbotapi.BotAPI
    userSvc     
    recordSvc
    handlers    map[string]func(tgbotapi.Update) error
}

