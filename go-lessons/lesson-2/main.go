package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int    `json:"id"`
	PostTitle string `json:"post_title"`
	Body      string `json:"body"`
}

var posts []Post

func main() {
	r := gin.Default()
	r.GET("/greet/:name", func(c *gin.Context) {
		name := c.Param("name")

		c.String(http.StatusOK, "Привет, %s", name)
	})

	r.Run()
}
