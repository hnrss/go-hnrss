package main

import (
	"time"
)

// <http://cyber.harvard.edu/rss/rss.html>
type RSS struct {
	XMLName       string    `xml:"rss"`
	Version       string    `xml:"version,attr"`
	Title         string    `xml:"channel>title"`
	Link          string    `xml:"channel>link"`
	Description   string    `xml:"channel>description"`
	Webmaster     string    `xml:"channel>webMaster"`
	Docs          string    `xml:"channel>docs"`
	Generator     string    `xml:"channel>generator"`
	LastBuildDate string    `xml:"channel>lastBuildDate"`
	Items         []RSSItem `xml:"channel>item"`
}

type RSSPermalink struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

type RSSItem struct {
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	Link        string       `xml:"link"`
	Author      string       `xml:"author"`
	Comments    string       `xml:"comments"`
	Published   string       `xml:"pubDate"`
	Permalink   RSSPermalink `xml:"guid"`
}

func NewRSS(results *AlgoliaResponse, op *OutputParams) *RSS {
	rss := RSS{
		Version:       "2.0",
		Title:         op.Title,
		Link:          op.Link,
		Description:   "Hacker News RSS",
		Webmaster:     "https://github.com/edavis/go-hnrss/issues",
		Docs:          "https://edavis.github.io/go-hnrss/",
		Generator:     "https://github.com/edavis/go-hnrss",
		LastBuildDate: Timestamp("rss", time.Now().UTC()),
	}

	for _, hit := range results.Hits {
		item := RSSItem{
			Title:       hit.GetTitle(),
			Link:        hit.GetURL(op.LinkTo),
			Description: hit.GetDescription(),
			Author:      hit.Author,
			Comments:    hit.GetPermalink(),
			Published:   Timestamp("rss", hit.GetCreatedAt()),
			Permalink:   RSSPermalink{hit.GetPermalink(), "false"},
		}
		rss.Items = append(rss.Items, item)
	}

	return &rss
}
