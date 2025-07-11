package tmpl

import (
	"html/template"
	"time"
)

// FormatDate formats a time.Time value into a readable string with
// day, abbreviated month, year, and 24-hour time (e.g., "02 Jan 2006 15:04").
func FormatDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}

// FormatChip returns a CSS class name corresponding to the given action type.
// It maps "CREATE", "UPDATE", and "DELETE" strings to predefined chip class names
// used for styling status indicators in the UI.
func FormatChip(s string) template.CSS {
	switch s {
	case "CREATE":
		return "chip-create"
	case "UPDATE":
		return "chip-update"
	case "DELETE":
		return "chip-delete"
	}

	return ""
}
