package model

import "time"

type Link struct {
	ShortName   string
	LongName    string
	Visits      int
	LastUpdated time.Time
}
