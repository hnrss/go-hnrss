package main

import (
	"time"
)

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

// Docs:
// - RSS: http://cyber.harvard.edu/rss/rss.html
// - Atom: https://validator.w3.org/feed/docs/atom.html
// - JSONFeed: https://jsonfeed.org/version/1

// Feel free to open an issue if you discover something wonky going on
// with any of these three formats (esp. Atom and JSONFeed):
// <https://github.com/edavis/go-hnrss/issues/new>

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
	Relationship string `xml:"rel,attr,omitempty"`
	Type         string `xml:"type,attr,omitempty"`
}

type RSSPermalink struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

type RSSItem struct {
	Title       CDATA        `xml:"title"`
	Description CDATA        `xml:"description"`
	Link        string       `xml:"link"`
	Author      string       `xml:"dc:creator"`
	Comments    string       `xml:"comments"`
	Published   string       `xml:"pubDate"`
	Permalink   RSSPermalink `xml:"guid"`
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
			Title:       CDATA{hit.GetTitle()},
			Link:        hit.GetURL(op.LinkTo),
			Description: CDATA{hit.GetDescription()},
			Author:      hit.Author,
			Comments:    hit.GetPermalink(),
			Published:   Timestamp("rss", hit.GetCreatedAt()),
			Permalink:   RSSPermalink{hit.GetPermalink(), "false"},
		}
		rss.Items = append(rss.Items, item)
	}

	return &rss
}

type Atom struct {
	XMLName string     `xml:"feed"`
	NS      string     `xml:"xmlns,attr"`
	ID      string     `xml:"id"`
	Title   string     `xml:"title"`
	Updated string     `xml:"updated"`
	Links   []AtomLink `xml:"link"`
	Entries []AtomEntry
}

type AtomEntry struct {
	XMLName   string      `xml:"entry"`
	Title     CDATA       `xml:"title"`
	Links     []AtomLink  `xml:"link"`
	Author    AtomPerson  `xml:"author"`
	Content   AtomContent `xml:"content"`
	Updated   string      `xml:"updated"`
	Published string      `xml:"published"`
	ID        string      `xml:"id"`
}

type AtomPerson struct {
	Name string `xml:"name"`
}

type AtomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",cdata"`
}

func NewAtom(results *AlgoliaSearchResponse, op *OutputParams) *Atom {
	atom := Atom{
		NS:      NSAtom,
		ID:      "https://hnrss.org" + op.SelfLink,
		Title:   op.Title,
		Updated: Timestamp("atom", time.Now().UTC()),
		Links: []AtomLink{
			AtomLink{"https://hnrss.org" + op.SelfLink, "self", "application/atom+xml"},
		},
	}

	for _, hit := range results.Hits {
		if op.TopLevel && !hit.isTopLevelComment() {
			continue
		}

		entry := AtomEntry{
			ID:        hit.GetPermalink(),
			Title:     CDATA{hit.GetTitle()},
			Updated:   Timestamp("atom", hit.GetCreatedAt()),
			Published: Timestamp("atom", hit.GetCreatedAt()),
			Links: []AtomLink{
				AtomLink{hit.GetURL(op.LinkTo), "alternate", ""},
			},
			Author:  AtomPerson{hit.Author},
			Content: AtomContent{"html", hit.GetDescription()},
		}
		atom.Entries = append(atom.Entries, entry)
	}

	return &atom
}

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
