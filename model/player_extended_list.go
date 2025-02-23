package model

import "strings"

type PlayerExtendedList []PlayerExtended

func (p PlayerExtendedList) String() string {
	var sb strings.Builder
	sb.Grow(len(p) * 512)
	for _, player := range p {
		sb.WriteString(player.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (p PlayerExtendedList) ToPlayerList() PlayerList {
	players := make([]Player, 0, len(p))
	for _, player := range p {
		players = append(players, Player{
			Nickname: player.Nickname,
			IP:       player.IP,
			Text:     player.Text,
		})
	}
	return players
}

func (p PlayerExtendedList) ToIPList() StringList {
	ips := make(StringList, 0, len(p))
	for _, player := range p {
		ips = append(ips, player.IP)
	}
	return ips
}
