package model

type IPText struct {
	IP   string `json:"ip"`
	Text string `json:"text"`
}

type IPTextList []IPText
