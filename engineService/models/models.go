package models

type Engine struct {
	ID     int32   `json:"ID"` // TODO: почему json тег с заглавных букв?
	Volume float32 `json:"volume"`
}
