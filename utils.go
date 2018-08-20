package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
	SiteURL      = "https://hnrss.org"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

func Timestamp(fmt string, input time.Time) string {
	switch fmt {
	case "rss":
		return input.Format(time.RFC1123Z)
	case "atom", "jsonfeed":
		return input.Format(time.RFC3339)
	default:
		return input.Format(time.RFC1123Z)
	}
}

func UTCNow() time.Time {
	return time.Now().UTC()
}

func ParseRequest(c *gin.Context) (*SearchParams, *OutputParams) {
	var (
		sp SearchParams
		op OutputParams
	)

	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	if strings.Contains(sp.Query, " OR ") {
		sp.Query = strings.Replace(sp.Query, " OR ", " ", -1)

		var q []string
		for _, f := range strings.Fields(sp.Query) {
			q = append(q, fmt.Sprintf("\"%s\"", f))
		}
		sp.Query = strings.Join(q, " ")
		sp.OptionalWords = strings.Join(q, " ")
	}

	err = c.ShouldBindQuery(&op)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}
	op.Format = c.GetString("format")
	op.SelfLink = SiteURL + c.Request.URL.String()

	return &sp, &op
}

func Generate(c *gin.Context, sp *SearchParams, op *OutputParams) {
	if op.Format == "" {
		op.Format = "rss"
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Header("X-Algolia-URL", algoliaSearchURL+sp.Values().Encode())

	switch op.Format {
	case "rss":
		rss := NewRSS(results, op)
		c.XML(http.StatusOK, rss)
	case "atom":
		atom := NewAtom(results, op)
		c.XML(http.StatusOK, atom)
	case "jsonfeed":
		jsonfeed := NewJSONFeed(results, op)
		c.JSON(http.StatusOK, jsonfeed)
	}
}
