package clock

import "time"

// Clock defines a function provider for obtaining current timestamps.
// This abstraction allows for easier unit testing by mocking the time source.
type Clock func() int64

// UnixMilliClock returns the current Unix time in milliseconds.
// It is the default implementation used across the logging system.
func UnixMilliClock() int64 {
	return time.Now().UnixMilli()
}
