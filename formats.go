package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	NS_DC   = "http://purl.org/dc/elements/1.1/"
	NS_ATOM = "http://www.w3.org/2005/Atom"
)

///////////////////////////////////////////////////////////////////////////

type RSS struct {
	XMLName    string `xml:"rss"`
	DublinCore string `xml:"xmlns:dc,attr"`
	Atom       string `xml:"xmlns:atom,attr"`
	Version    string `xml:"version,attr"`

	Title         string `xml:"channel>title"`
	Link          string `xml:"channel>link"`
	Description   string `xml:"channel>description"`
	Docs          string `xml:"channel>docs"`
	Generator     string `xml:"channel>generator"`
	LastBuildDate string `xml:"channel>lastBuildDate"`

	Items []RSSItem `xml:"channel>item"`
}

type RSSGuid struct {
	XMLName   string `xml:"guid"`
	Permalink string `xml:"isPermalink,attr"`
	Value     string `xml:",chardata"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	URL         string `xml:"link"`
	Creator     string `xml:"dc:creator"`
	Comments    string `xml:"comments"`
	UniqueID    RSSGuid
}

func OutputRSS(results AlgoliaResponse, c *gin.Context) {
	rss := RSS{
		Version:    "2.0",
		DublinCore: NS_DC,
		Atom:       NS_ATOM,
	}
	rss.Title = "Hello World from RSS"

	for _, hit := range results.Hits {
		rss.Items = append(rss.Items, RSSItem{
			Title:    hit.Title,
			URL:      hit.URL,
			UniqueID: RSSGuid{Permalink: "false", Value: hit.URL},
		})
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
