package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OutputParams struct {
	Title       string
	Link        string
	Description string `form:"description"`
	LinkTo      string `form:"link"`
	Format      string `form:"format"`
}

func (op OutputParams) Output(c *gin.Context, results *AlgoliaResponse) {
	fmt := op.Format
	if fmt == "" {
		fmt = "rss"
	}

	switch fmt {
	case "rss":
		op.RSS(c, results)
	}
}

func (op OutputParams) RSS(c *gin.Context, results *AlgoliaResponse) {
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

	c.XML(http.StatusOK, rss)
}

type SearchParams struct {
	Tags             string
	Query            string `form:"q"`
	Points           string `form:"points"`
	Comments         string `form:"comments"`
	SearchAttributes string `form:"search_attrs"`
	Count            string `form:"count"`
}

func (sp *SearchParams) numericFilters() string {
	var filters []string
	if sp.Points != "" {
		filters = append(filters, "points>="+sp.Points)
	}
	if sp.Comments != "" {
		filters = append(filters, "num_comments>="+sp.Comments)
	}
	return strings.Join(filters, ",")
}

// Encode transforms the search options into an Algolia search querystring
func (sp *SearchParams) Values() url.Values {
	params := make(url.Values)

	if sp.Query != "" {
		params.Set("query", sp.Query)
	}

	if f := sp.numericFilters(); f != "" {
		params.Set("numericFilters", f)
	}

	searchAttrs := sp.SearchAttributes
	if searchAttrs == "" {
		searchAttrs = "title"
	}
	if searchAttrs != "default" {
		params.Set("restrictSearchableAttributes", searchAttrs)
	}

	if sp.Count != "" {
		params.Set("hitsPerPage", sp.Count)
	}

	if sp.Tags != "" {
		params.Set("tags", sp.Tags)
	}

	return params
}

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/newest", Newest)
	r.GET("/frontpage", FrontPage)
	r.GET("/newcomments", NewComments)
	r.GET("/ask", Ask)
	r.GET("/show", Show)
	r.GET("/polls", Polls)
	r.GET("/jobs", Jobs)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	r.Run()
}
