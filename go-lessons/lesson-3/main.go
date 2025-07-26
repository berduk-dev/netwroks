package main

import (
	"lesson-3/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	postsHandler := handler.Handler{
		LastID: 0,
		Posts:  make(map[int]handler.Post),
	}

	r.POST("/posts", postsHandler.CreatePostHandler)
	r.GET("posts/:id", postsHandler.GetPostHandler)
	r.GET("posts", postsHandler.GetPostsHandler)
	r.DELETE("posts/:id", postsHandler.DeletePostHandler)
	r.PUT("posts/:id", postsHandler.UpdatePostHandler)

	r.Run()
}
