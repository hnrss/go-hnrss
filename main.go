package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	algoliaEndpoint = "https://hn.algolia.com/api/v1/search_by_date"
)

///////////////////////////////////////////////////////////////////////////

func Newest(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	opts := make(url.Values)
	opts.Add("tags", "(story,poll)")

	if query := qs.Get("q"); query != "" {
		opts.Add("query", query)
	}

	var nf []string
	if points := qs.Get("points"); points != "" {
		nf = append(nf, "points>="+points)
	}
	if comments := qs.Get("comments"); comments != "" {
		nf = append(nf, "num_comments>="+comments)
	}
	if len(nf) > 0 {
		opts.Add("numericFilters", strings.Join(nf, ","))
	}

	io.WriteString(w, algoliaEndpoint+"?"+opts.Encode()+"\n")
}

///////////////////////////////////////////////////////////////////////////

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/newest", Newest)
	log.Fatal(http.ListenAndServe(":9001", r))
}
