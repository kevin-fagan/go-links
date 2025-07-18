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
		return "bg-green-100 text-green-800"
	case "UPDATE":
		return "bg-yellow-100 text-yellow-800"
	case "DELETE":
		return "bg-red-100 text-red-800"
	}

	return ""
}
