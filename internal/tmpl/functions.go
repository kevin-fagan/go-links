package tmpl

import "time"

// FormatDate formats a time.Time value into a readable string with
// day, abbreviated month, year, and 24-hour time (e.g., "02 Jan 2006 15:04").
func FormatDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}
