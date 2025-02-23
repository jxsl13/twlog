package match

import (
	"regexp"
	"strconv"
)

var (
	// 0: full 1: ID 2: IP
	ddnetJoinRegex = regexp.MustCompile(`(?i)player has entered the game\. ClientID=([\d]+) addr=[^\d]{0,2}([\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3})[^\d]{0,2}`)

	// 0: full 1: ID 2: IP 3: port 4: version 5: name 6: clan 7: country
	playerzCatchJoinRegex = regexp.MustCompile(`(?i)id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+):([\d]+) version=(\d+) name='(.{0,20})' clan='(.{0,16})' country=([-\d]+)$`)

	// 0: full 1: ID 2: IP
	playerVanillaJoinRegex = regexp.MustCompile(`(?i)player is ready\. ClientID=([\d]+) addr=[^\d]{0,2}([\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3})[^\d]{0,2}`)
)

func Join(line string) (id int, ip string, ok bool) {
	var (
		joinIDStr string
		joinIP    string
	)
	if matches := ddnetJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else if matches := playerzCatchJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else if matches := playerVanillaJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else {
		return -1, "", false
	}

	joinId, err := strconv.Atoi(joinIDStr)
	if err != nil {
		return -1, "", false
	}

	return joinId, joinIP, true
}
