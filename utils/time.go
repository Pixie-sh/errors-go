package utils

import "time"

// TimeToIso convert time value to predefined ISO format
func TimeToIso(t time.Time) string {
	return t.Format(time.RFC3339)
}
