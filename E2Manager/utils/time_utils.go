package utils

import "time"

func ElapsedTime(startTime time.Time) float64 {
	return float64(time.Since(startTime)) / float64(time.Millisecond)
}
