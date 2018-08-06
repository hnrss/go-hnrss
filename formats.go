package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	NS_DC   = "http://purl.org/dc/elements/1.1/"
	NS_ATOM = "http://www.w3.org/2005/Atom"
)

///////////////////////////////////////////////////////////////////////////

type RSS struct {
	XMLName    string `xml:"rss"`
	DublinCore string `xml:"xmlns:dc,attr"`
	Version    string `xml:"version,attr"`
	Debug      string `xml:",comment"`

	Title         string `xml:"channel>title"`
	Link          string `xml:"channel>link"`
	Description   string `xml:"channel>description"`
	Docs          string `xml:"channel>docs"`
	Generator     string `xml:"channel>generator"`
	LastBuildDate string `xml:"channel>lastBuildDate"`

	Items []RSSItem `xml:"channel>item"`
}

type RSSPermalink struct {
	Value       string `xml:",chardata"`
	IsPermalink string `xml:"isPermalink,attr"`
}

type RSSItem struct {
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	URL         string       `xml:"link"`
	Creator     string       `xml:"dc:creator"`
	Comments    string       `xml:"comments"`
	Published   string       `xml:"pubDate"`
	Permalink   RSSPermalink `xml:"guid"`
}

func OutputRSS(results AlgoliaResponse, c *gin.Context) {
	rss := RSS{
		Version:    "2.0",
		DublinCore: NS_DC,
		Debug:      c.GetString("request_url"),

		Title:         c.GetString("output_title"),
		Link:          c.GetString("output_link"),
		Description:   "Hacker News RSS",
		Docs:          "https://edavis.github.io/hnrss/",
		Generator:     "https://github.com/edavis/hnrss",
		LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
	}

	for _, hit := range results.Hits {
		permalink := "https://news.ycombinator.com/item?id=" + hit.ObjectID
		created_at, _ := time.Parse("2006-01-02T15:04:05.000Z", hit.CreatedAt)

		item := RSSItem{
			Title:     hit.Title,
			Creator:   hit.Author,
			Comments:  permalink,
			Published: created_at.Format(time.RFC1123Z),
			Permalink: RSSPermalink{permalink, "false"},
		}

		switch c.GetString("resp_link") {
		case "url":
			item.URL = hit.URL
		case "comments":
			item.URL = permalink
		}

		rss.Items = append(rss.Items, item)
	}

	c.XML(http.StatusOK, rss)
}

///////////////////////////////////////////////////////////////////////////

type Atom struct {
	Title string `xml:"feed>title"`
}

func OutputAtom(c *gin.Context) {
	atom := Atom{}
	atom.Title = "hello world from Atom"
	c.XML(http.StatusOK, atom)
}

///////////////////////////////////////////////////////////////////////////

type JsonFeed map[string]interface{}

func OutputJsonFeed(c *gin.Context) {
	jsonfeed := make(JsonFeed)
	c.JSON(http.StatusOK, jsonfeed)
}
