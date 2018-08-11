package main

import "time"

func Timestamp(fmt string, input time.Time) string {
	switch fmt {
	case "rss":
		return input.Format(time.RFC1123Z)
	case "jsonfeed":
		return input.Format(time.RFC3339)
	default:
		return input.Format(time.RFC1123Z)
	}
}
