package chatmapper

import (
	"gopkg.in/Iwark/spreadsheet.v2"

	"github.com/arykalin/metaver-topos/users"
)

type Links struct {
	TrackName           string
	VKGroup             string
	TrackChat           string
	FindTeamTopic       string
	MentorRoleChat      string
	LeaderRoleChat      string
	VolunteerRoleChat   string
	ParticipantRoleChat string
}

const (
	TrackRow = "Трек"
)

type ChatMap = map[string]Links

type chatmapper struct {
	ChatMap ChatMap
}

type ChatMapper interface {
	GetMap() ChatMap
}

func (c *chatmapper) GetMap() ChatMap {
	return c.ChatMap
}

func NewChatMapper(sheet *spreadsheet.Sheet) ChatMapper {
	/*
		Mentor           = "Наставник"
		TrackLeader      = "Ведущий трека"
		Participant      = "Участник"
		Volunteer        = "Волонтер"
	*/
	var (
		MentorRoleChat      string
		LeaderRoleChat      string
		VolunteerRoleChat   string
		ParticipantRoleChat string
	)

	chatMap := make(ChatMap)
	for i := range sheet.Rows {
		if sheet.Rows[i][0].Value == "В каком качестве вы хотите участвовать в марафоне?" {
			switch sheet.Rows[i][1].Value {
			case users.Mentor:
				MentorRoleChat = sheet.Rows[i][6].Value
			case users.TrackLeader:
				LeaderRoleChat = sheet.Rows[i][6].Value
			case users.Participant:
				ParticipantRoleChat = sheet.Rows[i][6].Value
			case users.Volunteer:
				VolunteerRoleChat = sheet.Rows[i][6].Value
			}
		}
	}

	for i := range sheet.Rows {
		if i == 0 {
			//skip 0 row
			continue
		}
		if sheet.Rows[i][0].Value != TrackRow {
			continue
		}
		//log.Printf("chat map values: %s", sheet.Values[i])
		track := sheet.Rows[i][1].Value
		vkGroup := sheet.Rows[i][2].Value
		trackChat := sheet.Rows[i][3].Value
		findTeamTopic := sheet.Rows[i][4].Value

		chatMap[track] = Links{
			TrackName:           track,
			VKGroup:             vkGroup,
			TrackChat:           trackChat,
			FindTeamTopic:       findTeamTopic,
			MentorRoleChat:      MentorRoleChat,
			LeaderRoleChat:      LeaderRoleChat,
			VolunteerRoleChat:   VolunteerRoleChat,
			ParticipantRoleChat: ParticipantRoleChat,
		}
	}

	return &chatmapper{
		ChatMap: chatMap,
	}
}
