package date

import (
	"time"
)

const (
	ApiDateTimeFormat  = "01-02-2006T11:06:39:000Z"
	ApiTimeStampFormat = "2006-01-02 15:04:05"
)

func GetTimeNow() string {
	return time.Now().UTC().Format(ApiDateTimeFormat)
}
func GetTimeNowDB() string {
	return time.Now().UTC().Format(ApiTimeStampFormat)
	//return time.Now().UTC().Format(time.RFC1123)
}
