package utils

import (
	"fmt"
	"time"
)

// ParseDate parses a date string into a time.Time object. The date string must be in YYYY-MM-DD format.
func ParseDate(dateStr string) (time.Time, error) {
	d, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error in parsing date %s: %w", dateStr, err)
	}
	return d, nil
}
