package date

import (
	"time"
)

const (
	ApiDateTimeFormat = "01-02-2006T11:06:39:000Z"
)

func GetTimeNow() string {
	return string(time.Now().UTC().Format(ApiDateTimeFormat))
}
