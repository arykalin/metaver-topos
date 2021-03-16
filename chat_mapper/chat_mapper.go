package chatmapper

import (
	"gopkg.in/Iwark/spreadsheet.v2"
)

type Links struct {
	TrackName   string
	VKGroup     string
	TrackChat   string
	VkTeamTopic string
	RoleChat    string
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
		vkTeamTopic := sheet.Rows[i][4].Value
		roleChat := sheet.Rows[i][6].Value

		chatMap[track] = Links{
			TrackName:   track,
			VKGroup:     vkGroup,
			TrackChat:   trackChat,
			VkTeamTopic: vkTeamTopic,
			RoleChat:    roleChat,
		}
	}
	chatMap["Среда как зона поиска, влияния и поддержки"] = Links{
		TrackName:   "Среда как зона поиска, влияния и поддержки",
		VKGroup:     "https://vk.com/kfest_edu_sreda",
		TrackChat:   "https://vk.me/join/AJQ1d8kQGhubwIIzQhVOKPK/",
		VkTeamTopic: "",
		RoleChat:    "",
	}
	chatMap["Среда как зона поиска, влияния и поддержки (трек для педагогов)"] = Links{
		TrackName:   "Среда как зона поиска, влияния и поддержки (трек для педагогов)",
		VKGroup:     "https://vk.com/kfest_edu_sreda",
		TrackChat:   "https://vk.me/join/AJQ1d8kQGhubwIIzQhVOKPK/",
		VkTeamTopic: "",
		RoleChat:    "",
	}
	chatMap["Открытые образовательные пространства: программа и сообщество (трек для взрослых)"] = Links{
		TrackName:   "Открытые образовательные пространства: программа и сообщество (трек для взрослых)",
		VKGroup:     "https://vk.com/kfest_open_edu_space",
		TrackChat:   "https://t.me/joinchat/g5ndt4A3FSdiZjli",
		VkTeamTopic: "https://vk.com/topic-203127412_47177535",
		RoleChat:    "",
	}
	return &chatmapper{
		ChatMap: chatMap,
	}
}
