package model

import "fmt"

type PlayerExtended struct {
	File     string `json:"file"`
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Text     string `json:"text"`
}

func (p PlayerExtended) String() string {
	return fmt.Sprintf("%s: id=%d ip=%s name=%s text=%s", p.File, p.ID, p.IP, p.Nickname, p.Text)
}
