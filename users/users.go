package users

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	userDataJsonFile = "user_data.json"
	Mentor           = "Наставник"
	TrackLeader      = "Ведущий трека"
	Participant      = "Участник"
	Volunteer        = "Волонтер"
)

const (
	UserTypeUnknonwn UserType = iota
	UserTypeMentor
	UserTypeVolunteer
	UserTypeParticipant
	UserTypeTrackLeader
)

type UserType int

type User struct {
	Email    string
	Track    string
	Name     string
	HaveTeam bool
	Type     UserType
}

type Users = map[string]User

type SheetConfig struct {
	TrackIdx    int
	MailIdx     int
	HaveTeam    bool
	UserTypeIdx *int
	Skip        int
}

type users struct {
	users Users
}

type UsersInt interface {
	AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig) (err error)
	GetUsers() Users
	DumpUsers() error
}

func (u *users) AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig) (err error) {
	for i := range sheet.Rows {
		if i < config.Skip {
			// skip
			continue
		}
		user := User{}
		if config.UserTypeIdx != nil {
			switch sheet.Rows[i][*config.UserTypeIdx].Value {
			case Mentor:
				user.Type = UserTypeMentor
			case TrackLeader:
				user.Type = UserTypeTrackLeader
			case Volunteer:
				user.Type = UserTypeVolunteer
			case Participant:
				user.Type = UserTypeParticipant
			}
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
