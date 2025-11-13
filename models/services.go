package models

type Services struct {
	ID     int    `json:"id"`
	Eng    string `json:"eng"`
	Title  string `json:"title"`
	Prices string `json:"prices"`
	Src    string `json:"src"`
	Text   string `json:"text"`
}
