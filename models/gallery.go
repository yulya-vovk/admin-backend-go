package models

type Gallery struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Hidden   bool   `json:"hidden"`
}
