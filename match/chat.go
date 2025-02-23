package match

import (
	"regexp"
	"strconv"
)

var (
	// id, nick, chat line
	chatLineRegexp = regexp.MustCompile(`chat: (\d+):-?\d+:(.+): (.+)`)
)

func Chat(line string) (id int, nick string, chat string, ok bool) {
	matches := chatLineRegexp.FindStringSubmatch(line)
	if len(matches) == 0 {
		return -1, "", "", false
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// must match, otherwise hte regex is wrong
		panic(err)
	}
	nick = matches[2]
	chat = matches[3]
	return id, nick, chat, true
}
