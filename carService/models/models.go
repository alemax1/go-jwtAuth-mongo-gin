package models

import _ "encoding/json"

type Car struct {
	ID       int32  `json:"ID"`
	Concern  string `json:"concern"`
	Model    string `jsom:"model"`
	Year     int32  `json:"year"`
	Used     bool   `json:"used"`
	EngineID int32  `json:"engineID"`
	Price    int32  `json:"price"`
}
