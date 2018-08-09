package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const algoliaURL = "https://hn.algolia.com/api/v1/search_by_date?"

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

func (op OutputParams) getURL(hit AlgoliaHit, permalink string) string {
	linkTo := op.LinkTo
	if linkTo == "" {
		linkTo = "url"
	}

	switch {
	case linkTo == "url" && hit.URL != "":
		return hit.URL
	case linkTo == "url" && hit.URL == "":
		return permalink
	case linkTo == "comments":
		return permalink
	default:
		return permalink
	}
}

func (op OutputParams) getDescription(hit AlgoliaHit) string {
	if hit.isComment() {
		return hit.CommentText
	} else {
		return "default description"
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
		LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
	}

	for _, hit := range results.Hits {
		permalink := "https://news.ycombinator.com/item?id=" + hit.ObjectID
		createdAt, _ := time.Parse("2006-01-02T15:04:05.000Z", hit.CreatedAt)
		item := RSSItem{
			Title:       hit.GetTitle(),
			Link:        op.getURL(hit, permalink),
			Description: op.getDescription(hit),
			Author:      hit.Author,
			Comments:    permalink,
			Published:   createdAt.Format(time.RFC1123Z),
			Permalink:   RSSPermalink{permalink, "false"},
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

type AlgoliaHit struct {
	ObjectID    string
	Title       string
	URL         string
	Author      string
	CreatedAt   string   `json:"created_at"`
	StoryTitle  string   `json:"story_title"`
	CommentText string   `json:"comment_text"`
	Tags        []string `json:"_tags"`
}

func (hit AlgoliaHit) isComment() bool {
	for _, tag := range hit.Tags {
		if tag == "comment" {
			return true
		}
	}
	return false
}

func (hit AlgoliaHit) GetTitle() string {
	if hit.isComment() {
		return fmt.Sprintf("New comment by %s in \"%s\"", hit.Author, hit.StoryTitle)
	} else {
		return hit.Title
	}
}

type AlgoliaResponse struct {
	Hits []AlgoliaHit
}

func GetResults(params url.Values) (*AlgoliaResponse, error) {
	resp, err := http.Get(algoliaURL + params.Encode())
	if err != nil {
		return nil, errors.New("Error getting search results from Algolia")
	}
	defer resp.Body.Close()

	var parsed AlgoliaResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&parsed)
	if err != nil {
		return nil, errors.New("Invalid JSON received from Algolia")
	}

	return &parsed, nil
}

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/newest", Newest)
	r.GET("/ask", Ask)
	r.GET("/show", Show)
	r.GET("/newcomments", NewComments)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	r.Run()
}
