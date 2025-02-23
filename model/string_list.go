package model

import "strings"

type StringList []string

func (s StringList) String() string {
	var sb strings.Builder
	sb.Grow(len(s) * 64)
	for _, str := range s {
		sb.WriteString(str)
		sb.WriteByte('\n')
	}
	return sb.String()
}
