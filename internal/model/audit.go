package model

import "time"

type Audit struct {
	ShortURL  string
	LongURL   string
	Action    string
	Timestamp time.Time
}
