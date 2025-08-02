package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type Handler struct {
	db *pgx.Conn
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

func NewHandler(db *pgx.Conn) Handler {
	return Handler{
		db,
	}
}

func (h *Handler) CreatePost(c *gin.Context) {
	var post Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	_, err = h.db.Exec(c, "INSERT INTO posts (title, body) VALUES ($1, $2)", post.Title, post.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Что-то пошло не так!")
		log.Println("error insert into posts ", err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPostById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("ошибка в strconv.Atoi:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невалидный id"})
		return
	}

	var post Post

	err = h.db.QueryRow(c, "SELECT id, title, body, created_at FROM posts WHERE id = $1", id).Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.CreatedAt,
	)

	if err != nil {
		log.Println("ошибка при выполнении запроса:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	rows, err := h.db.Query(c, "SELECT id, title, body, created_at FROM posts")
	if err != nil {
		log.Println("ошибка запроса:", err)
		c.JSON(http.StatusInternalServerError, "Ошибка сервера")
		return
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt)
		if err != nil {
			log.Println("ошибка сканирования строки:", err)
			c.JSON(http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, "Постов нет")
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id")
		log.Println("error in DeletePostHandler in strconv.Atoi: ", err)
		return
	}

	err = h.db.QueryRow(c, "SELECT id, title, body, created_at FROM posts WHERE id = $1", id).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, "Такого поста нет!")
		return
	}

	_, err = h.db.Exec(c, "DELETE FROM posts WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Что-то пошло не так!")
		log.Println("error post delete", err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Пост №%d удалён!", id))
}

func (h *Handler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id!")
		log.Println("error in UpdatePostHandler in strconv.Atoi: ", err)
		return
	}

	err = h.db.QueryRow(c, "SELECT id, title, body, created_at FROM posts WHERE id = $1", id).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, "Такого поста нет!")
		return
	}

	var updatePostRequest UpdatePostRequest
	err = c.BindJSON(&updatePostRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос!")
		return
	}

	_, err = h.db.Exec(c, "UPDATE posts SET title = $1, body = $2", updatePostRequest.Title, updatePostRequest.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Что-то пошло не так!")
		log.Println("error post update ", err)
		return
	}

	err = h.db.QueryRow(c, "SELECT id, title, body, created_at FROM posts WHERE id = $1", id).Scan(
		&updatePostRequest.Title,
		&updatePostRequest.Body,
	)
	if err != nil {
		log.Println("Ошибка при изменении поста: ", err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Пост №%d изменён!", id))
}
