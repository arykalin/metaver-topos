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

const (
	track = "Среда как зона поиска, влияния и поддержки (трек для педагогов)"
)

type Mailer interface {
	SendGreeting(user users.User, info chatmapper.Links) error
}

const (
	tmpltMentorWithoutTeam      = "mailer/mails/template_mentor_without_team.html"
	tmpltMentorWithTeam         = "mailer/mails/template_mentor_with_team.html"
	tmpltParticipantWithTeam    = "mailer/mails/template_participant_with_team.html"
	tmpltParticipantWithoutTeam = "mailer/mails/template_participant_without_team.html"
	tmpltVolunteerWithoutTeam   = "mailer/mails/template_volunteer_without_team.html"
)

func (m mailer) SendGreeting(user users.User, info chatmapper.Links) (err error) {
	if info.TrackName == "" {
		return fmt.Errorf("track name is empty")
	}

	body, subj, err := m.makeBodyAndSubj(user, info)
	if err != nil {
		return err
	}

	if body == "" {
		return fmt.Errorf("body is empty")
	}

	if subj == "" {
		return fmt.Errorf("subject is empty")
	}

	err = m.sendMail(user, subj, body)
	if err != nil {
		return err
	}
	return nil
}

func (m mailer) sendMail(user users.User, subj string, body string) (err error) {
	//if debug address is set send mail to it
	var toEmail string
	if m.DebugAddress != "" {
		toEmail = m.DebugAddress
	} else {
		toEmail = user.Email
	}

	gm := gomail.NewMessage()
	gm.SetHeader("From", m.User)
	gm.SetHeader("To", toEmail)
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
	return err
}

func (m mailer) makeBodyAndSubj(user users.User, info chatmapper.Links) (body string, subj string, err error) {
	switch user.Type {
	case users.UserTypeUnknonwn:
		return body, subj, fmt.Errorf("user type unknown")
	case users.UserTypeTrackLeader:
		if user.HaveTeam {
			subj = fmt.Sprintf("Регистрация команды в Краефест - трек \"%s\"", user.Track)
			body, err = m.ParseTemplate(tmpltParticipantWithTeam, info)
			if err != nil {
				return body, subj, err
			}
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltParticipantWithTeam)
		}
	case users.UserTypeMentor:
		if user.HaveTeam {
			body, err = m.ParseTemplate(tmpltMentorWithTeam, info)
			if err != nil {
				return body, subj, err
			}
			subj = fmt.Sprintf("Регистрация наставника с командой в Краефест  - трек \"%s\"", user.Track)
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltMentorWithoutTeam)
		}
		if !user.HaveTeam {
			body, err = m.ParseTemplate(tmpltMentorWithoutTeam, info)
			if err != nil {
				return body, subj, err
			}
			subj = fmt.Sprintf("Регистрация наставника в Краефест  - трек \"%s\"", user.Track)
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltMentorWithoutTeam)
		}
	case users.UserTypeParticipant:
		if user.HaveTeam {
			subj = fmt.Sprintf("Регистрация команды в Краефест - трек \"%s\"", user.Track)
			body, err = m.ParseTemplate(tmpltParticipantWithTeam, info)
			if err != nil {
				return body, subj, err
			}
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltParticipantWithTeam)
		}
		if !user.HaveTeam {
			subj = fmt.Sprintf("Регистрация участника в Краефест - трек \"%s\"", user.Track)
			body, err = m.ParseTemplate(tmpltParticipantWithoutTeam, info)
			if err != nil {
				return body, subj, err
			}
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltParticipantWithoutTeam)
		}
	case users.UserTypeVolunteer:
		if !user.HaveTeam {
			subj = fmt.Sprintf("Регистрация волонтера в Краефест - трек \"%s\"", user.Track)
			body, err = m.ParseTemplate(tmpltVolunteerWithoutTeam, info)
			if err != nil {
				return body, subj, err
			}
			m.logger.Debugw("user template", "user", user.Email, "template", tmpltVolunteerWithoutTeam)
		}
	default:
		return body, subj, fmt.Errorf("can not determine user type for template")
	}
	return body, subj, err
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
