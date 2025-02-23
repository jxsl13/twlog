package model

import (
	"fmt"
	"strings"
)

type Player struct {
	Nickname string `json:"nickname"`
	IP       string `json:"ip"`
	Text     string `json:"text"`
}

func (p Player) String() string {
	return fmt.Sprintf("<{%s}> %s: %s", p.IP, p.Nickname, p.Text)
}

type PlayerList []Player

func (p PlayerList) String() string {
	var sb strings.Builder
	sb.Grow(len(p) * 256)
	for _, player := range p {
		sb.WriteString(player.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}
