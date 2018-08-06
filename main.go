package main

import (
	"encoding/json"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

const algoliaUrl = "https://hn.algolia.com/api/v1/search_by_date?"

type AlgoliaResponse struct {
	Hits []struct {
		Title string
		URL   string
	}
}

func Prepare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract out all the available query parameters
		c.Set("req_query", c.Query("q"))
		c.Set("req_id", c.Query("id"))
		c.Set("req_points", c.Query("points"))
		c.Set("req_comments", c.Query("comments"))
		c.Set("req_search_attrs", c.DefaultQuery("search_attrs", "title"))
		c.Set("req_count", c.DefaultQuery("count", "25"))
		c.Set("resp_link", c.DefaultQuery("link", "url"))
		c.Set("resp_description", c.DefaultQuery("description", "on"))
		c.Set("resp_format", c.DefaultQuery("format", "rss"))

		// Build the query string for Algolia
		params := make(url.Values)

		// Attach tags
		params.Set("tags", c.GetString("req_tags"))

		// Attach query
		// TODO(ejd): try query[] for OR?
		if query := c.GetString("req_query"); query != "" {
			params.Set("query", query)
		}

		// Attach points and/or comments filter
		var filters []string
		if points := c.GetString("req_points"); points != "" {
			filters = append(filters, "points>="+points)
		}
		if comments := c.GetString("req_comments"); comments != "" {
			filters = append(filters, "num_comments>="+comments)
		}
		if len(filters) > 0 {
			params.Set("numericFilters", strings.Join(filters, ","))
		}

		// Attach search attributes
		if search_attrs := c.GetString("req_search_attrs"); search_attrs != "" && search_attrs != "default" {
			params.Set("restrictSearchableAttributes", search_attrs)
		}

		// Attach count
		// TODO(ejd): cap this at 100
		if count := c.GetString("req_count"); count != "" {
			params.Set("hitsPerPage", count)
		}

		// TODO(ejd): put together a smarter HTTP client
		resp, err := http.Get(algoliaUrl + params.Encode())
		if err != nil {
			c.String(http.StatusBadGateway, "Error getting search results from Algolia")
		}
		defer resp.Body.Close()

		var parsed AlgoliaResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&parsed); err != nil {
			c.String(http.StatusBadGateway, "Invalid JSON received from Algolia")
		}

		switch c.GetString("resp_format") {
		case "json":
			OutputJsonFeed(c)
		case "atom":
			OutputAtom(c)
		case "rss":
			OutputRSS(parsed, c)
		default:
			c.String(http.StatusBadRequest, "Format must be one of: 'json', 'atom', or 'rss'")
		}

		c.Next()
	}
}

///////////////////////////////////////////////////////////////////////////

func Newest(c *gin.Context) {
	c.Set("req_tags", "(story,poll)")
}

func Frontpage(c *gin.Context) {
	c.Set("req_tags", "front_page")
}

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/newest", Newest, Prepare())
	r.GET("/frontpage", Frontpage, Prepare())

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	r.Run()
}
