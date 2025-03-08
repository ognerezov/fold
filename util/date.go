package util

import "time"

func IsTimePast(t time.Time, delta int) bool {
	now := time.Now()
	return t.Sub(now) < time.Second*time.Duration(delta)
}
