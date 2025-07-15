package logs

import "time"

type Logs struct {
	// Short URL associated with the log.
	// This will be populated when changes are made to a link
	ShortURL string
	// Long URL associated with the log.
	// This will be populated when changes are made to a link
	LongURL string
	// Tag associated with the log.
	// This will be populated when changes are made to a tag
	Tag string
	// The action (create, update, delete) taken on a link or a tag
	Action string
	// The IP of the client. OAuth is currently not supported so this will do for now
	ClientIP string
	// Represents when the log was created
	Timestamp time.Time
}
