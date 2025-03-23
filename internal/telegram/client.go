package telegram

import (
	"p2p-stats/internal/domain/repositories"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Client struct {
    Bot             *tgbotapi.BotAPI
    recordRepo      repositories.RecordRepo
    userRepo        repositories.UserRepository
    handlers        map[string]func(tgbotapi.Update) error
}

func NewClient(token string, recordRepo repositories.RecordRepo, userRepo repositories.UserRepository) *Client {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Fatal().Err(err).Msg("Error creating tg bot")
    }
    c := &Client{
        Bot:      bot,
        recordRepo: recordRepo,
        userRepo: userRepo,
        handlers: make(map[string]func(tgbotapi.Update) error),
    }
    c.registerHandlers()
    return c
}

func (c *Client) registerHandlers() {
    c.handlers["/start"] = func(update tgbotapi.Update) error {
        // Create a new user in the database
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!")
        _, err := c.Bot.Send(msg)
        return err
    }
    c.handlers["/sell"] = func(update tgbotapi.Update) error {
        // Create a new sell record
        return nil
    }
    c.handlers["/buy"] = func(update tgbotapi.Update) error {
        // Create a new sell record
        return nil
    }
}

func (c *Client) ProcessUpdate(update tgbotapi.Update) {
    if update.Message != nil && update.Message.Text != "" {
        cmd := strings.Split(update.Message.Text, " ")[0]
        if handler, ok := c.handlers[cmd]; ok {
            if err := handler(update); err != nil {
                log.Fatal().Err(err).Msg("Error processing command")
            }
        }
    }
}
