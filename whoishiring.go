package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func HiringCommon(c *gin.Context, query string) {
	params := make(url.Values)
	params.Set("query", fmt.Sprintf("\"%s\"", query))
	params.Set("tags", "story,author_whoishiring")
	params.Set("hitsPerPage", "1")

	results, err := GetResults(params)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	if len(results.Hits) < 1 {
		e := errors.New("No whoishiring stories found")
		c.Error(e)
		c.String(http.StatusBadGateway, e.Error())
		return
	}

	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	sp.Filters = "parent_id=" + results.Hits[0].ObjectID
	sp.SearchAttributes = "default"
	op.Title = results.Hits[0].Title
	op.Link = "https://news.ycombinator.com/item?id=" + results.Hits[0].ObjectID

	Generate(c, &sp, &op)
}

func SeekingEmployees(c *gin.Context) {
	HiringCommon(c, "Ask HN: Who is hiring?")
}

func SeekingEmployers(c *gin.Context) {
	HiringCommon(c, "Ask HN: Who wants to be hired?")
}

func SeekingFreelance(c *gin.Context) {
	HiringCommon(c, "Ask HN: Freelancer? Seeking freelancer?")
}
