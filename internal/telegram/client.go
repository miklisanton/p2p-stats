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
    toDelete        map[int64][]int
}

func NewClient(token string, recordRepo repositories.RecordRepo, userRepo repositories.UserRepository) *Client {
    bot, err := tgbotapi.NewBotAPI(token)
    bot.Debug = true
    if err != nil {
        log.Fatal().Err(err).Msg("Error creating tg bot")
    }
    c := &Client{
        Bot:      bot,
        recordRepo: recordRepo,
        userRepo: userRepo,
        handlers: make(map[string]func(tgbotapi.Update) error),
        toDelete: make(map[int64][]int),
    }
    c.registerHandlers()
    return c
}



func (c *Client) registerHandlers() {
    c.handlers["/start"] = c.HandleStart
    c.handlers["/sell"] = c.HandleAddRecord
    c.handlers["/buy"] =  c.HandleAddRecord
    c.handlers["/list"] =  c.HandleListRecords
    c.handlers["del"] = c.HandleDeleteRecord
    c.handlers["exp"] = c.HandleExportRecords
}

func (c *Client) Start() {
    c.SetCommands()
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := c.Bot.GetUpdatesChan(u)
    for update := range updates {
        c.ProcessUpdate(update)
    }
}

func (c *Client) SetCommands() {
    commands := []tgbotapi.BotCommand{
        {Command: "start", Description: "Start using the bot"},
        {Command: "list", Description: "List orders"},
    }
    cfg := tgbotapi.NewSetMyCommands(commands...)
    _, err := c.Bot.Request(cfg)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to set commands")
    }
}

func (c *Client) ProcessUpdate(update tgbotapi.Update) {
    if update.Message != nil && update.Message.Text != "" {
        cmd := "/" + update.Message.Command()
        log.Info().Str("cmd", cmd).Msg("Command")
        if handler, ok := c.handlers[cmd]; ok {
            c.DeleteMessages(update.Message.Chat.ID)
            if err := handler(update); err != nil {
                log.Error().Err(err).Msg("Error processing command")
            }
        } else {
            c.SendMessage(update.Message.Chat.ID, "Unknown command")
        }
    }

    if update.CallbackQuery != nil {
        cmd := strings.Split(update.CallbackQuery.Data, "|")[0]
        if handler, ok := c.handlers[cmd]; ok {
            if err := handler(update); err != nil {
                log.Error().Err(err).Msg("Error processing callback")
            }
        } else {
            log.Info().Str("cmd", cmd).Msg("No handler for callback")
            c.SendMessage(update.CallbackQuery.Message.Chat.ID, "Internal error")
        }
    }
}

func (c *Client) SendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := c.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Str("msg", message).Msg("Failed to send message")
	}
}

func (c *Client) SendMessageNDelete(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	m, err := c.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Str("msg", message).Msg("Failed to send message")
	}
	c.DeleteAppend(chatID, m.MessageID)
}

// Delete messages from toDelete list and clear it
func (c *Client) DeleteMessages(chatID int64) {
	for _, id := range c.toDelete[chatID] {
		msg := tgbotapi.NewDeleteMessage(chatID, id)
		_, err := c.Bot.Send(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete message")
		}
	}
	c.toDelete[chatID] = nil
}

// Adds message to delete list
func (c *Client) DeleteAppend(chatID int64, msgID int) {
	if c.toDelete[chatID] == nil {
		c.toDelete[chatID] = make([]int, 0)
	}

	c.toDelete[chatID] = append(c.toDelete[chatID], msgID)
}
