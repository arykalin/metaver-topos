package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	chatmapper "github.com/arykalin/metaver-topos/chat_mapper"
	"github.com/arykalin/metaver-topos/mailer"
	"github.com/arykalin/metaver-topos/users"
)

type SentFile struct {
	Sent map[string]bool `json:"mail_sent,omitempty"`
}

type notifier struct {
	logger   *zap.SugaredLogger
	mailer   mailer.Mailer
	sentFile string
}

type Notifier interface {
	Notify(chatMap chatmapper.ChatMap, users users.Users) error
}

func (n notifier) Notify(chatMap chatmapper.ChatMap, users users.Users) error {
	sentData := SentFile{}
	sentData.Sent = make(map[string]bool)
	byteValue, err := ioutil.ReadFile(n.sentFile)
	if err != nil {
		return err
	}
	//make backup
	backupFile := fmt.Sprintf("%s-%d", n.sentFile, time.Now().Unix())
	err = ioutil.WriteFile(backupFile, byteValue, 0644)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &sentData)
	if err != nil {
		return err
	}

	for mail, user := range users {
		sentRecord := fmt.Sprintf("%s-%s", mail, user.Track)
		if sent, ok := sentData.Sent[sentRecord]; ok && sent {
			n.logger.Debugw("message already sent. skip", "mail", mail)
			continue
		}
		if trackInfo, ok := chatMap[user.Track]; ok {
			n.logger.Debugf("sending links to user %s from track: %s. track info: %+v\n", mail, user.Track, trackInfo)
			err = n.mailer.SendGreeting(user, trackInfo)
			if err != nil {
				n.logger.Errorw("sending mail error", "mail", mail, "track", user.Track, "type", user.Type.Name(), "err", err)
				continue
			}
			sentData.Sent[sentRecord] = true
		} else {
			n.logger.Debugf("user not found in track map mail: %s track: %s", mail, user.Track)
		}
	}

	file, err := json.MarshalIndent(sentData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(n.sentFile, file, 0644) //nolint:gosec
	return err
}

func NewNotifier(logger *zap.SugaredLogger, mailer mailer.Mailer, sentFile string) Notifier {
	return &notifier{
		logger:   logger,
		mailer:   mailer,
		sentFile: sentFile,
	}
}
