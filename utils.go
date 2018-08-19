package main

import "time"

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

func Timestamp(fmt string, input time.Time) string {
	switch fmt {
	case "rss":
		return input.Format(time.RFC1123Z)
	case "atom", "jsonfeed":
		return input.Format(time.RFC3339)
	default:
		return input.Format(time.RFC1123Z)
	}
}

func UTCNow() time.Time {
	return time.Now().UTC()
}
