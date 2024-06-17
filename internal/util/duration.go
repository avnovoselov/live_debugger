package util

import "time"

// IntToMillisecond - return duration millisecond
func IntToMillisecond(duration int) time.Duration {
	return time.Duration(duration) * time.Millisecond
}
