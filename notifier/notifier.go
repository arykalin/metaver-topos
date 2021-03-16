package notifier

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"go.uber.org/zap"

	chatmapper "github.com/arykalin/metaver-topos/chat_mapper"
	"github.com/arykalin/metaver-topos/mailer"
	"github.com/arykalin/metaver-topos/users"
)

const sentJsonFile = "sent.json"

type SentFile struct {
	Sent map[string]bool `json:"mail_sent,omitempty"`
}

type notifier struct {
	logger *zap.SugaredLogger
	mailer mailer.Mailer
}

type Notifier interface {
	Notify(chatMap chatmapper.ChatMap, users users.Users) error
}

func (n notifier) Notify(chatMap chatmapper.ChatMap, users users.Users) error {
	sentData := SentFile{}
	sentData.Sent = make(map[string]bool)
	jsonData, err := os.Open(sentJsonFile)
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonData)
	err = json.Unmarshal(byteValue, &sentData)
	if err != nil {
		return err
	}
	err = jsonData.Close()
	if err != nil {
		return err
	}

	for mail, user := range users {
		if sent, ok := sentData.Sent[mail]; ok && sent {
			n.logger.Debugw("message already sent. skip", "mail", mail)
			continue
		}
		if trackInfo, ok := chatMap[user.Track]; ok {
			n.logger.Debugf("sending links to user %s from track: %s. track info: %+v\n", mail, user.Track, trackInfo)
			err = n.mailer.SendGreeting(user, trackInfo)
			if err != nil {
				return err
			}
			sentData.Sent[mail] = true
		} else {
			n.logger.Debugf("user not found in track map mail: %s track: %s", mail, user.Track)
		}
	}

	file, err := json.MarshalIndent(sentData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(sentJsonFile, file, 0644) //nolint:gosec
	return err
}

func NewNotifier(logger *zap.SugaredLogger, mailer mailer.Mailer) Notifier {
	return &notifier{
		logger: logger,
		mailer: mailer,
	}
}
