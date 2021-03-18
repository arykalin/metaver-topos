package telegram

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendMessage(chatID int64, text string) error {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	msg := tgbotapi.NewMessage(chatID, text)
	//msg.ReplyToMessageID = update.Message.MessageID

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
