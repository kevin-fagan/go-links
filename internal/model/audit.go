package model

import "time"

type Action string

var (
	Create Action = "CREATE"
	Delete Action = "DELETE"
	Update Action = "UPDATE"
)

type Audit struct {
	ShortURL  string
	LongURL   string
	Action    Action
	Timestamp time.Time
}
