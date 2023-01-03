package timestamp

import (
	"time"
)

func StringToInt64(s string) int64 {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return tm.Unix()
}

func StringToInt64Micro(s string) int64 {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return tm.UnixNano()
}

func SecToDateTimeString() string {
	i := time.Now().Unix()
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func ShortTimeString() string {
	tm := time.Unix(time.Now().Unix(), 0)
	return tm.Format("01-02 15:04")
}

func GetDateString() string {
	return time.Now().Format("2006-01-02")
}

func NowTimeString() string {
	return time.Now().Format("2006-01-02_15:04")
}
