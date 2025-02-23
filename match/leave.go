package match

import (
	"regexp"
	"strconv"
)

var (
	// 0: full 1: ID 2: IP 3: reason
	playerLeftRegex = regexp.MustCompile(`id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+) reason='(.*)'$`)
)

func Leave(line string) (id int, ok bool) {
	matches := playerLeftRegex.FindStringSubmatch(line)
	if len(matches) == 0 {
		return -1, false
	}

	idStr := matches[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, false
	}

	return id, true
}
