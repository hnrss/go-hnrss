package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func runner(c *gin.Context, sp SearchParams, op OutputParams) {
	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.Header("X-Algolia-URL", algoliaURL+sp.Values().Encode())

	op.Output(c, results)
}

func Newest(c *gin.Context) {
	var (
		sp SearchParams
		op OutputParams
	)

	sp.Tags = "(story,poll)"
	op.Title = "Hacker News: Newest"
	op.Link = "https://news.ycombinator.com/newest"

	runner(c, sp, op)
}

func NewComments(c *gin.Context) {
	var (
		sp SearchParams
		op OutputParams
	)

	sp.Tags = "comment"
	sp.SearchAttributes = "default"
	op.Title = "Hacker News: New Comments"
	op.Link = "https://news.ycombinator.com/newcomments"

	runner(c, sp, op)
}
