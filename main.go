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
	Format      string
	SelfLink    string
	TopLevel    bool
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

func SetFormat(fmt string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("format", fmt)
		c.Next()
	}
}

func registerEndpoint(r *gin.Engine, url string, fn gin.HandlerFunc) {
	r.GET(url, SetFormat("rss"), fn)
	r.GET(url+".rss", SetFormat("rss"), fn)
	r.GET(url+".jsonfeed", SetFormat("jsonfeed"), fn)
	r.GET(url+".atom", SetFormat("atom"), fn)
}

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	registerEndpoint(r, "/newest", Newest)
	registerEndpoint(r, "/frontpage", Frontpage)
	registerEndpoint(r, "/newcomments", Newcomments)
	registerEndpoint(r, "/ask", AskHN)
	registerEndpoint(r, "/show", ShowHN)
	registerEndpoint(r, "/polls", Polls)
	registerEndpoint(r, "/jobs", Jobs)
	registerEndpoint(r, "/user", UserAll)
	registerEndpoint(r, "/threads", UserThreads)
	registerEndpoint(r, "/submitted", UserSubmitted)
	registerEndpoint(r, "/item", Item)
	registerEndpoint(r, "/whoishiring/jobs", SeekingEmployees)
	registerEndpoint(r, "/whoishiring/hired", SeekingEmployers)
	registerEndpoint(r, "/whoishiring/freelance", SeekingFreelance)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	r.Run()
}
