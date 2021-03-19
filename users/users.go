package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	userDataJsonFile = "user_data.json"
	Mentor           = "Наставник"
	TrackLeader      = "Ведущий трека"
	Participant      = "Участник"
	Volunteer        = "Волонтер"
	Unknown          = "Unknonwn"
)

const (
	UserTypeUnknonwn UserType = iota
	UserTypeMentor
	UserTypeVolunteer
	UserTypeParticipant
	UserTypeTrackLeader
)

type UserType int

func (u UserType) Name() string {
	switch u {
	case UserTypeUnknonwn:
		return Unknown
	case UserTypeMentor:
		return Mentor
	case UserTypeVolunteer:
		return Volunteer
	case UserTypeParticipant:
		return Participant
	case UserTypeTrackLeader:
		return TrackLeader
	}
	return Unknown
}

type FormType int

const (
	FormTypeNone FormType = iota
	FormTypeWithTeam
	FormTypeNoTeam
)

type User struct {
	Email       string
	LeaderEmail string
	Track       string
	Name        string
	HaveTeam    bool
	Type        UserType
}

type Users = map[string]User

type SheetConfig struct {
	TrackIdx      int
	MailIdx       int
	LeaderMailIdx int
	HaveTeam      bool
	UserTypeIdx   *int
	Skip          int
}

type users struct {
	users Users
}

type UsersInt interface {
	AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig, formType FormType) (err error)
	GetUsers() Users
	DumpUsers() error
}

func (u *users) AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig, formType FormType) (err error) {
	for i := range sheet.Rows {
		if i < config.Skip {
			// skip
			continue
		}
		user := User{}
		switch formType {
		case FormTypeNoTeam:
			switch sheet.Rows[i][*config.UserTypeIdx].Value {
			case Mentor:
				user.Type = UserTypeMentor
			case TrackLeader:
				user.Type = UserTypeTrackLeader
			case Volunteer:
				user.Type = UserTypeVolunteer
			case Participant:
				user.Type = UserTypeParticipant
			default:
				user.Type = UserTypeUnknonwn
			}
		case FormTypeWithTeam:
			user.Type = UserTypeMentor
		default:
			return fmt.Errorf("form type unknonwn")
		}

		var track string
		if len(sheet.Rows[i]) > config.TrackIdx {
			track = sheet.Rows[i][config.TrackIdx].Value
		}
		user.Track = track
		var mail string
		if len(sheet.Rows[i]) > config.MailIdx {
			mail = sheet.Rows[i][config.MailIdx].Value
		}

		if formType == FormTypeWithTeam {
			var leaderMail string
			if len(sheet.Rows[i]) > config.LeaderMailIdx {
				leaderMail = sheet.Rows[i][config.LeaderMailIdx].Value
			}

			if leaderMail != "" {
				user.Email = leaderMail
				user.HaveTeam = config.HaveTeam
				user.Type = UserTypeTrackLeader
				u.users[leaderMail] = user
			}
		}

		user.Email = mail
		user.HaveTeam = config.HaveTeam
		u.users[mail] = user
	}
	return err
}

func (u users) GetUsers() Users {
	return u.users
}

func (u users) DumpUsers() error {
	file, err := json.MarshalIndent(u.users, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(userDataJsonFile, file, 0644) //nolint:gosec
	return err
}
func NewUsers() UsersInt {
	makeUsers := make(Users)
	return &users{
		users: makeUsers,
	}
}
