package links

import "time"

type Link struct {
	ShortURL    string
	LongURL     string
	Visits      int
	LastUpdated time.Time
}
