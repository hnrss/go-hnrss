package main

import (
	"time"
)

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
)

// <http://cyber.harvard.edu/rss/rss.html>
type RSS struct {
	XMLName       string    `xml:"rss"`
	Version       string    `xml:"version,attr"`
	NSDublinCore  string    `xml:"xmlns:dc,attr"`
	NSAtom        string    `xml:"xmlns:atom,attr"`
	Title         string    `xml:"channel>title"`
	Link          string    `xml:"channel>link"`
	Description   string    `xml:"channel>description"`
	Docs          string    `xml:"channel>docs"`
	Generator     string    `xml:"channel>generator"`
	LastBuildDate string    `xml:"channel>lastBuildDate"`
	AtomLink      AtomLink  `xml:"channel>atom:link"`
	Items         []RSSItem `xml:"channel>item"`
}

type AtomLink struct {
	Reference    string `xml:"href,attr"`
	Relationship string `xml:"rel,attr"`
	Type         string `xml:"type,attr"`
}

type RSSPermalink struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

type RSSDescription struct {
	Value string `xml:",cdata"`
}

type RSSItem struct {
	Title       string         `xml:"title"`
	Description RSSDescription `xml:"description"`
	Link        string         `xml:"link"`
	Author      string         `xml:"dc:creator"`
	Comments    string         `xml:"comments"`
	Published   string         `xml:"pubDate"`
	Permalink   RSSPermalink   `xml:"guid"`
}

func NewRSS(results *AlgoliaSearchResponse, op *OutputParams) *RSS {
	rss := RSS{
		Version:       "2.0",
		NSAtom:        NSAtom,
		NSDublinCore:  NSDublinCore,
		Title:         op.Title,
		Link:          op.Link,
		Description:   "Hacker News RSS",
		Docs:          "https://edavis.github.io/go-hnrss/",
		Generator:     "https://github.com/edavis/go-hnrss",
		LastBuildDate: Timestamp("rss", time.Now().UTC()),
		AtomLink:      AtomLink{op.SelfLink, "self", "application/rss+xml"},
	}

	for _, hit := range results.Hits {
		if op.TopLevel && !hit.isTopLevelComment() {
			continue
		}

		item := RSSItem{
			Title:       hit.GetTitle(),
			Link:        hit.GetURL(op.LinkTo),
			Description: RSSDescription{hit.GetDescription()},
			Author:      hit.Author,
			Comments:    hit.GetPermalink(),
			Published:   Timestamp("rss", hit.GetCreatedAt()),
			Permalink:   RSSPermalink{hit.GetPermalink(), "false"},
		}
		rss.Items = append(rss.Items, item)
	}

	return &rss
}

// ----------------------------------------------------------------------

// <https://jsonfeed.org/version/1>
type JSONFeed struct {
	Version     string         `json:"version"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Link        string         `json:"home_page_url"`
	Items       []JSONFeedItem `json:"items"`
}

type JSONFeedItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ContentHTML string `json:"content_html"`
	URL         string `json:"url"`
	ExternalURL string `json:"external_url"`
	Published   string `json:"date_published"`
	Author      string `json:"author"`
}

func NewJSONFeed(results *AlgoliaSearchResponse, op *OutputParams) *JSONFeed {
	jf := JSONFeed{
		Version:     "https://jsonfeed.org/version/1",
		Title:       op.Title,
		Link:        op.Link,
		Description: "Hacker News RSS",
	}
	for _, hit := range results.Hits {
		if op.TopLevel && !hit.isTopLevelComment() {
			continue
		}

		item := JSONFeedItem{
			ID:          hit.GetPermalink(),
			Title:       hit.GetTitle(),
			ContentHTML: hit.GetDescription(),
			URL:         hit.GetURL(op.LinkTo),
			ExternalURL: hit.GetPermalink(),
			Published:   Timestamp("jsonfeed", hit.GetCreatedAt()),
			Author:      hit.Author,
		}
		jf.Items = append(jf.Items, item)
	}
	return &jf
}
