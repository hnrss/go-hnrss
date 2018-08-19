package main

import (
	"fmt"
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
	OptionalWords    string
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

	if sp.OptionalWords != "" {
		params.Set("query", sp.Query)
		params.Set("optionalWords", sp.OptionalWords)
	} else if sp.Query != "" {
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
