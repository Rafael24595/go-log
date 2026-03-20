package clock

import "time"

type Clock func() int64

func UnixMilliClock() int64 {
	return time.Now().UnixMilli()
}
