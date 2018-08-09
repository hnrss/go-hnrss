package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	hackerNewsItemID = "https://news.ycombinator.com/item?id="
	algoliaURL       = "https://hn.algolia.com/api/v1/search_by_date?"
)

type AlgoliaResponse struct {
	Hits []AlgoliaHit
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

func (hit AlgoliaHit) GetPermalink() string {
	return hackerNewsItemID + hit.ObjectID
}

func (hit AlgoliaHit) GetURL(linkTo string) string {
	if linkTo == "" {
		linkTo = "url"
	}

	switch {
	case linkTo == "url" && hit.URL != "":
		return hit.URL
	default:
		return hit.GetPermalink()
	}
}

func (hit AlgoliaHit) GetDescription() string {
	if hit.isComment() {
		return hit.CommentText
	} else {
		return "" // TODO(ejd)
	}
}

func (hit AlgoliaHit) GetCreatedAt() time.Time {
	rv, err := time.Parse("2006-01-02T15:04:05.000Z", hit.CreatedAt)
	if err != nil {
		return time.Now().UTC()
	} else {
		return rv
	}
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
