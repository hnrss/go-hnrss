package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func registerEndpoint(r *gin.Engine, url string, fn gin.HandlerFunc) {
	r.GET(url, SetFormat("rss"), fn)
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
	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-agent: *\nDisallow:\n")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	var addr []string
	if port := os.Getenv("PORT"); port != "" {
		addr = append(addr, "127.0.0.1:"+port)
	}
	r.Run(addr...)
}
