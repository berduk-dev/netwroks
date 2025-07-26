package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Handler struct {
	LastID int
	Posts  map[int]Post
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

func (h *Handler) CreatePostHandler(c *gin.Context) {
	var post Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	h.LastID++
	post.ID = h.LastID
	h.Posts[h.LastID] = post
	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPostHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id")
		log.Println("error in GetPostHandler in strconv.Atoi: ", err)
		return
	}

	post, ok := h.Posts[id]
	if !ok {
		c.JSON(http.StatusNotFound, "Такого поста нет!")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPostsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.Posts)
}

func (h *Handler) DeletePostHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id")
		log.Println("error in DeletePostHandler in strconv.Atoi: ", err)
		return
	}

	_, ok := h.Posts[id]
	if !ok {
		c.JSON(http.StatusNotFound, "Такого поста нет!")
		return
	}

	delete(h.Posts, id)

	c.JSON(http.StatusOK, fmt.Sprintf("Пост №%d удалён", id))
}

func (h *Handler) UpdatePostHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id!")
		log.Println("error in UpdatePostHandler in strconv.Atoi: ", err)
		return
	}

	var updatePostRequest UpdatePostRequest
	err = c.BindJSON((&updatePostRequest))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос!")
		return
	}

	post, ok := h.Posts[id]
	if !ok {
		c.JSON(http.StatusNotFound, "Такого поста нет!")
		return
	}

	if updatePostRequest.Body != nil {
		post.Body = *updatePostRequest.Body
	}
	if updatePostRequest.Title != nil {
		post.Title = *updatePostRequest.Title
	}
	h.Posts[id] = post

	c.JSON(http.StatusOK, h.Posts[id])
}
