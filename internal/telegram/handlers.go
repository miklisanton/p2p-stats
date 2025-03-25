package telegram

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"p2p-stats/internal/domain/entities"
	"p2p-stats/internal/utils"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const WelcomeMessage = "🇬🇧 Welcome to P2P Stats Bot!\n" +
"Commands:\n" +
"*Buy 100 USDT for 90 EUR*\n" +
"`/buy 100 90 EUR`\n" +
"*Sell 100 USDT for 9800 RUB*\n" +
"`/sell 100 9800 RUB`\n" +
"*All Trades*\n" +
"`/list`\n" +
"*From date*\n" +
"`/list 01-01-2024`\n" +
"*In range*\n" +
"`/list 01-01-2024 31-12-2024`\n\n" +
"🇷🇺 Добро пожаловать в P2P Stats Bot!\n" +
"Команды:\n" +
"*Покупка 100 USDT за 90 EUR*\n" +
"`/buy 100 90 EUR`\n" +
"*Продажа 100 USDT за 9800 RUB*\n" +
"`/sell 100 9800 RUB`\n" +
"*Все сделки*\n" +
"`/list`\n" +
"*С даты*\n" +
"`/list 01-01-2024`\n" +
"*За период*\n" +
"`/list 01-01-2024 31-12-2024`"

func (c *Client)HandleStart(update tgbotapi.Update) error {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, WelcomeMessage)
    msg.ParseMode = "markdown"
    if _, err := c.Bot.Send(msg); err != nil {
        return err
    }

    if err := c.userRepo.Create(update.Message.Chat.ID); err != nil {
        return err
    }
    return nil
}

func (c *Client)HandleAddRecord(update tgbotapi.Update) error {
    // Create a new record
    r, err := utils.ParseRecord(entities.RecordType(strings.ToUpper(update.Message.Command())), update.Message.CommandArguments())
    if err != nil {
        c.SendMessage(update.Message.Chat.ID, err.Error())
        return err
    }
    r.UserID = update.Message.Chat.ID
    valR, err := entities.NewValidatedRecord(r)
    if err != nil {
        c.SendMessage(update.Message.Chat.ID, err.Error())
        return err
    }
    if err := c.recordRepo.Create(valR); err != nil {
        return err
    }
    return nil
}

func(c *Client)HandleListRecords(update tgbotapi.Update) error {
    c.DeleteAppend(update.Message.Chat.ID, update.Message.MessageID)
    args := update.Message.CommandArguments()
    from, to, err := utils.ParseList(args)
    if err != nil {
        c.SendMessage(update.Message.Chat.ID, err.Error())
        return err
    }
    
    records, err := c.recordRepo.GetByUserID(update.Message.Chat.ID, from, to)
    if err != nil {
        return err
    }

    if len(records) == 0 {
        c.SendMessageNDelete(update.Message.Chat.ID, "No records found for given timeframe")
        return nil
    }

    if args == "" {
        args = "all"
    }

    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    for _, r := range records {
        row := tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("❌" + r.String(), "del|" + r.Id.String() + "|" + args),
        )
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

    }
    expRow := tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("📎Export csv", "exp|" + args),
    )
    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, expRow)
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Delete or export records")
    msg.ReplyMarkup = keyboard
    m, err := c.Bot.Send(msg)
    if err != nil {
        return err
    }
    c.DeleteAppend(update.Message.Chat.ID, m.MessageID)
    return nil
}

func (c *Client)HandleDeleteRecord(update tgbotapi.Update) error {
    parts := strings.Split(update.CallbackQuery.Data, "|")
    if len(parts) != 3 {
        return nil
    }
    id := parts[1]
    args := parts[2]
    if args == "all" {
        args = ""
    }
    if err := c.recordRepo.Delete(id); err != nil {
        return err
    }
    c.DeleteMessages(update.CallbackQuery.Message.Chat.ID)
    // Reuse the list handler
    update.Message = &tgbotapi.Message{
        Chat: &tgbotapi.Chat{
            ID: update.CallbackQuery.Message.Chat.ID,
        },
        Text: "/list " + args,
    }
    return c.HandleListRecords(update)
}

func (c *Client)HandleExportRecords(update tgbotapi.Update) error {
    timeframe := strings.Split(update.CallbackQuery.Data, "|")[1]
    if timeframe == "all" {
        timeframe = ""
    }
    from, to, err := utils.ParseList(timeframe)
    if err != nil {
        return err
    }
    log.Info().Msgf("Exporting records from %v to %v", from, to)
    records, err := c.recordRepo.GetByUserID(update.CallbackQuery.Message.Chat.ID, from, to)
    if err != nil {
        return err
    }

    var buf bytes.Buffer

    // Add UTF-8 BOM (optional, ensures proper encoding in some applications like Excel)
    buf.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM

    // Create a CSV writer
    writer := csv.NewWriter(&buf)
    writer.Comma = ',' // Default delimiter, can be changed to ';' or other

    csvData := [][]string{
        {"Order No.", "Type", "Fiat Amount", "Currency", "Price", "Coin amount", "Cryptocurrency", "Time"},
    }
    for _, r := range records {
        csvData = append(csvData, []string{
            r.Id.String(),
            string(r.Type),
            fmt.Sprintf("%.2f", r.FiatAmount),
            r.FiatCurrency,
            fmt.Sprintf("%.2f", r.FiatAmount / r.USDTAmount),
            fmt.Sprintf("%.2f", r.USDTAmount),
            "USDT",
            r.CreatedAt.Format("2006-01-02 15:04:05"),
        })

    }
   if err := writer.WriteAll(csvData); err != nil {
       return err
   }
   writer.Flush()

   name := "all"
   if to == nil && from != nil {
       now := time.Now()
       to = &now
   }

   if to != nil && from != nil {
       name = from.Format("02-01-2006") + "_" + to.Format("02-01-2006")
   }

   fileBytes := tgbotapi.FileBytes{
       Name:  "orders_" + name + ".csv",
       Bytes: buf.Bytes(),
   }

   doc := tgbotapi.NewDocument(update.CallbackQuery.Message.Chat.ID, fileBytes)
   _, err = c.Bot.Send(doc)
   if err != nil {
       return err
   }
    
   c.DeleteMessages(update.CallbackQuery.Message.Chat.ID)
   return nil
}


