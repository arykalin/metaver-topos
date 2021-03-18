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
	SMTPPort     int
	DebugAddress string
	CCAddress    string
}

type Mailer interface {
	SendGreeting(user users.User, info chatmapper.Links) error
}

func (m mailer) SendGreeting(user users.User, info chatmapper.Links) (err error) {
	var (
		body string
		subj string
	)

	//if debug address is set send mail to it
	if m.DebugAddress != "" {
		user.Email = m.DebugAddress
	}

	if info.TrackName == "" {
		return fmt.Errorf("track name is empty")
	}

	switch user.Type {
	case users.UserTypeUnknonwn:
		return fmt.Errorf("user type unknown")
	case users.UserTypeMentor:
		body, err = m.ParseTemplate("mailer/mails/template_mentor_without_team.html", info)
		if err != nil {
			return err
		}
		subj = fmt.Sprintf("Регистрация участника в КраеФест - трек \"%s\"", user.Track)
		m.logger.Debugw("user template", "user", user.Email, "body", body)
	case users.UserTypeParticipant:
		if user.HaveTeam {
			body, err = m.ParseTemplate("mailer/mails/template_participant_with_team.html", info)
			if err != nil {
				return err
			}
			subj = fmt.Sprintf("Регистрация команды в Краефест - %s", user.Track)
			m.logger.Debugw("user template", "user", user.Email, "body", body)
		}
		if !user.HaveTeam {
			subj = fmt.Sprintf("Регистрация в Краефест - %s", user.Track)
			body, err = m.ParseTemplate("mailer/mails/template_participant_without_team.html", info)
			if err != nil {
				return err
			}
			m.logger.Debugw("user template", "user", user.Email, "body", body)
		}
	case users.UserTypeVolunteer:

	default:
		return fmt.Errorf("can not determine user type for template")
	}

	if body == "" {
		return fmt.Errorf("body is empty")
	}

	if subj == "" {
		return fmt.Errorf("subject is empty")
	}

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
	port int,
	debugAddress string,
	ccAddress string,
) Mailer {
	return &mailer{
		logger:       logger,
		User:         user,
		Password:     password,
		SMTPHost:     host,
		SMTPPort:     port,
		DebugAddress: debugAddress,
		CCAddress:    ccAddress,
	}
}
