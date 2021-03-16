package users

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/Iwark/spreadsheet.v2"
)

const userDataJsonFile = "user_data.json"

type User struct {
	Email    string
	Track    string
	Name     string
	HaveTeam bool
	IsMentor bool
}

type Users = map[string]User

type users struct {
	users Users
}

type UsersInt interface {
	AddUsers(sheet *spreadsheet.Sheet, trackIdx int, mailIdx int, haveTeam bool, mentorIdx int) (err error)
	GetUsers() Users
	DumpUsers() error
}

func (u *users) AddUsers(sheet *spreadsheet.Sheet, trackIdx int, mailIdx int, haveTeam bool, mentorIdx int) (err error) {
	for i := range sheet.Rows {
		if i == 0 {
			//skip header
			continue
		}
		user := User{}
		if mentorIdx != 0 {
			mentorField := sheet.Rows[i][mentorIdx].Value
			if mentorField == "Наставник" {
				user.IsMentor = true
			}
		}
		var track string
		if len(sheet.Rows[i]) > trackIdx {
			track = sheet.Rows[i][trackIdx].Value
		}
		user.Track = track
		var mail string
		if len(sheet.Rows[i]) > mailIdx {
			mail = sheet.Rows[i][mailIdx].Value
		}
		user.Email = mail
		user.HaveTeam = haveTeam
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
