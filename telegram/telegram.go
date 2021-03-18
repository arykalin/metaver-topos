package telegram

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type teleLog struct {
	chatID int64
	bot    *tgbotapi.BotAPI
}

type TeleLog interface {
	SendMessage(text string) error
}

func (t *teleLog) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(t.chatID, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func NewTelegramLog(chatID int64, token string) TeleLog {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &teleLog{
		chatID: chatID,
		bot:    bot,
	}
}
