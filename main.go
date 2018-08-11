package main

// TODO(ejd): create enums for the different formats

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

type OutputParams struct {
	Title       string
	Link        string
	Description string `form:"description"`
	LinkTo      string `form:"link"`
	Format      string `form:"format"`
}

type SearchParams struct {
	Tags             string
	Query            string `form:"q"`
	Points           string `form:"points"`
	ID               string `form:"id"`
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
		params.Set("query", fmt.Sprintf("\"%s\"", sp.Query))
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
	r.GET("/frontpage", Frontpage)
	r.GET("/newcomments", Newcomments)
	r.GET("/ask", AskHN)
	r.GET("/show", ShowHN)
	r.GET("/polls", Polls)
	r.GET("/jobs", Jobs)
	r.GET("/user", UserAll)
	r.GET("/threads", UserThreads)
	r.GET("/submitted", UserSubmitted)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	r.Run()
}
