package main

import (
	"log"
	"os"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("tgsectok")

	bot, err := telego.NewBot(botToken, telego.WithDefaultLogger(false, true))
	if err != nil {
		log.Panic(err)
	}

	defer bot.StopLongPolling()

	updates, _ := bot.UpdatesViaLongPolling(nil)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton("Удалить").WithCallbackData("RemoveMessage"),
		),
	)

	bh, _ := th.NewBotHandler(bot, updates)

	keywords := []string{"сделать", "заехать", "заказать", "записать", "купить", "позвонить", "пришить", "приклеить", "забрать",
		"to do", "to drive", "to order", "to write", "to buy", "to call", "to sew on", "to glue", "to pick up"}

	// handle any message
	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		if stringHasKeywordAsPrefix(strings.ToLower(message.Text), keywords) {
			msg := tu.Message(
				message.Chat.ChatID(),
				message.Text,
			).WithReplyMarkup(keyboard).WithProtectContent() // Multiple `with` method

			bot.SendMessage(msg)

			bot.DeleteMessage(&telego.DeleteMessageParams{
				ChatID:    message.Chat.ChatID(),
				MessageID: message.MessageID,
			})
		}
	})

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		bot.DeleteMessage(
			&telego.DeleteMessageParams{
				ChatID: telego.ChatID{
					ID: query.Message.GetChat().ID,
				},
				MessageID: query.Message.GetMessageID(),
			},
		)
	})

	defer bh.Stop()
	bh.Start()
}

func stringHasKeywordAsPrefix(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.HasPrefix(s, kw) {
			return true
		}
	}

	return false
}
