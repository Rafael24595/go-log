package format

import "time"

// FormatMilliseconds converts a Unix millisecond timestamp into a
// human-readable date and time string using the layout "YYYY-MM-DD HH:MM:SS".
// Returns "N/A" if the timestamp is zero.
func FormatMilliseconds(timestamp int64) string {
	if timestamp == 0 {
		return "N/A"
	}

	seconds := timestamp / 1000
	time := time.Unix(seconds, 0)

	return time.Format("2006-01-02 15:04:05")
}

// FormatMillisecondsCompact converts a Unix millisecond timestamp into a
// filesystem-friendly string using the layout "YYYYMMDD_HHMMSS".
// This is ideal for naming log files or snapshots.
// Returns "N/A" if the timestamp is zero.
func FormatMillisecondsCompact(timestamp int64) string {
	if timestamp == 0 {
		return "N/A"
	}
	seconds := timestamp / 1000
	time := time.Unix(seconds, 0)
	return time.Format("20060102_150405")
}
