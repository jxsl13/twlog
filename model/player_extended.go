package model

import (
	"fmt"

	"github.com/jxsl13/twlog/stringutils"
)

type PlayerExtended struct {
	File     string `json:"file"`
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Text     string `json:"text"`
}

func NewPlayerExtended(file, nickname string, id int, ip, text string) PlayerExtended {
	return PlayerExtended{
		File:     file,
		Nickname: stringutils.VisualizeInvisible(nickname),
		ID:       id,
		IP:       ip,
		Text:     stringutils.VisualizeInvisible(text),
	}
}

func (p PlayerExtended) String() string {
	return fmt.Sprintf("%s: id=%d ip=%s name=%s text=%s", p.File, p.ID, p.IP, p.Nickname, p.Text)
}
