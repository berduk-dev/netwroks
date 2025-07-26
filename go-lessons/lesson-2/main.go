package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int    `json:"id"`
	PostTitle string `json:"post_title"`
	Body      string `json:"body"`
}

func main() {
	r := gin.Default()
	r.GET("/greet/:name", func(c *gin.Context) {
		name := c.Param("name")

		c.String(http.StatusOK, "Привет, %s", name)
	})

	r.GET("/calc", func(c *gin.Context) {
		xStr := c.Query("x")
		yStr := c.Query("y")

		x, err := strconv.Atoi(xStr)
		if err != nil {
			log.Println("Ошибка в форматировании X:", err)
		}

		y, err := strconv.Atoi(yStr)
		if err != nil {
			log.Println("Ошибка в форматировании Y:", err)
		}

		sum := x + y
		c.JSON(http.StatusOK, gin.H{"sum": sum})
	})

	r.Run()
}
