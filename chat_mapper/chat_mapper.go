package chatmapper

import (
	"gopkg.in/Iwark/spreadsheet.v2"
)

type Links struct {
	TrackName     string
	VKGroup       string
	TrackChat     string
	FindTeamTopic string
	RoleChat      string
}

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
	chatMap := make(ChatMap)
	for i := range sheet.Rows {
		if i == 0 {
			//skip 0 row
			continue
		}

		//log.Printf("chat map values: %s", sheet.Values[i])
		track := sheet.Rows[i][1].Value
		vkGroup := sheet.Rows[i][2].Value
		trackChat := sheet.Rows[i][3].Value
		findTeamTopic := sheet.Rows[i][4].Value
		roleChat := sheet.Rows[i][6].Value

		chatMap[track] = Links{
			TrackName:     track,
			VKGroup:       vkGroup,
			TrackChat:     trackChat,
			FindTeamTopic: findTeamTopic,
			RoleChat:      roleChat,
		}
	}
	return &chatmapper{
		ChatMap: chatMap,
	}
}
