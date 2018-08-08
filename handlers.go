package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Newest(c *gin.Context) {
	var sp SearchParams
	sp.Tags = "(story,poll)"
	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	var op OutputParams
	err = c.ShouldBindQuery(&op)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}
	op.Title = "Hacker News: Newest"
	op.Link = "https://news.ycombinator.com/newest"
	op.Output(c, results)
}

func NewComments(c *gin.Context) {
	var sp SearchParams
	sp.Tags = "comment"
	sp.SearchAttributes = "default"
	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.Header("X-Algolia-URL", algoliaURL+sp.Values().Encode())

	var op OutputParams
	err = c.ShouldBindQuery(&op)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}
	op.Title = "Hacker News: New Comments"
	op.Link = "https://news.ycombinator.com/newcomments"
	op.Output(c, results)
}
