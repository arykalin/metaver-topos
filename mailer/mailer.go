package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"go.uber.org/zap"

	chatmapper "github.com/arykalin/metaver-topos/chat_mapper"
	"github.com/arykalin/metaver-topos/users"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	logger       *zap.SugaredLogger
	User         string
	Password     string
	SMTPHost     string
	DebugAddress string
	CCAddress    string
}

type Mailer interface {
	SendGreeting(user users.User, info chatmapper.Links) error
}

func (m mailer) SendGreeting(user users.User, info chatmapper.Links) (err error) {
	var body string

	//if debug address is set send mail to it
	if m.DebugAddress != "" {
		user.Email = m.DebugAddress
	}

	if info.TrackName == "" {
		return fmt.Errorf("track name is empty")
	}

	if user.IsMentor {
		body, err = m.ParseTemplate("mailer/mails/template_mentor_without_team.html", info)
		if err != nil {
			return err
		}
		m.logger.Debugw("user template", "user", user.Email, "body", body)
	}

	if !user.IsMentor {
		if user.HaveTeam {
			body, err = m.ParseTemplate("mailer/mails/template_with_team.html", info)
			if err != nil {
				return err
			}
			m.logger.Debugw("user template", "user", user.Email, "body", body)
		}
		if !user.HaveTeam {
			body, err = m.ParseTemplate("mailer/mails/template_without_team.html", info)
			if err != nil {
				return err
			}
			m.logger.Debugw("user template", "user", user.Email, "body", body)
		}
	}

	if body == "" {
		return fmt.Errorf("body is empty")
	}

	subj := fmt.Sprintf("Регистрация участника в КраеФест - трек \"%s\"", user.Track)
	gm := gomail.NewMessage()
	gm.SetHeader("From", m.User)
	gm.SetHeader("To", user.Email)
	if m.CCAddress != "" {
		gm.SetHeader("Cc", m.CCAddress)
	}
	gm.SetHeader("Subject", subj)
	gm.SetBody("text/html", body)

	d := gomail.NewDialer(m.SMTPHost, 587, m.User, m.Password)

	// Send the email.
	if err = d.DialAndSend(gm); err != nil {
		return err
	}
	return nil
}

func (m *mailer) ParseTemplate(templateFileName string, data interface{}) (body string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return body, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return body, err
	}
	return buf.String(), err
}

func NewMailer(
	logger *zap.SugaredLogger,
	user string,
	password string,
	host string,
	debugAddress string,
	ccAddress string,
) Mailer {
	return &mailer{
		logger:       logger,
		User:         user,
		Password:     password,
		SMTPHost:     host,
		DebugAddress: debugAddress,
		CCAddress:    ccAddress,
	}
}
