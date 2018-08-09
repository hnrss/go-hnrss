package main

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
