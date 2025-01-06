package utils

import "time"

func TimeNow() time.Time {
	return time.Now()
}

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

func TimeNowUnixMS() int64 {
	return time.Now().UnixMilli()
}
