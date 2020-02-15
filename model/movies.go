package model

type MovieInfo struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	DateLine     int64  `json:"dateline"`
	DateLineText string `json:"date_text"`
	HtmlUrl      string `json:"html_url"`
	VideoUrl     string `json:"video_url"`
	ThunderUrl   string `json:"thunder_url"`
	Finish       int    `json:"finish"`
	IsDown       int    `json:"is_down"`
}
