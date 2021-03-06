package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	chatmapper "github.com/arykalin/metaver-topos/chat_mapper"
	"github.com/arykalin/metaver-topos/mailer"
	"github.com/arykalin/metaver-topos/telegram"
	"github.com/arykalin/metaver-topos/users"
)

type SentFile struct {
	Sent map[string]bool `json:"mail_sent,omitempty"`
}

type notifier struct {
	logger   *zap.SugaredLogger
	mailer   mailer.Mailer
	sentFile string
	teleLog  telegram.TeleLog
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
				msg := fmt.Sprintf("sending mail error mail: %s track: %s type: %s sheet form: %s err: %s\n",
					mail, user.Track, user.Type.Name(), user.Form.Name(), err)
				n.logger.Error(msg)
				terr := n.teleLog.SendMessage(msg)
				if terr != nil {
					n.logger.Errorw("sending telegram message error", "err", err)
				}
				continue
			}
			msg := fmt.Sprintf("sent mail to user mail: %s\n track: %s\n type: %s\n haveTeam: %t\n track info:\n %+v\n sheet form: %s\n",
				mail, user.Track, user.Type.Name(), user.HaveTeam, trackInfo, user.Form.Name())
			terr := n.teleLog.SendMessage(msg)
			if terr != nil {
				n.logger.Errorw("sending telegram message error", "err", err)
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

func NewNotifier(
	logger *zap.SugaredLogger,
	mailer mailer.Mailer,
	sentFile string,
	teleLog telegram.TeleLog,
) Notifier {
	return &notifier{
		logger:   logger,
		teleLog:  teleLog,
		mailer:   mailer,
		sentFile: sentFile,
	}
}
