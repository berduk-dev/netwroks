package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"net/http"
)

const HostURL = "127.0.0.1:8080/"

type CreateLinkRequest struct {
	Link string `json:"link"`
}
type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) Handler {
	return Handler{
		db,
	}
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}
	var shortLink string
	// Проверка на наличие длинной ссылки в БД
	row := h.db.QueryRow(c, "SELECT short_link FROM links WHERE long_link = $1", req.Link)
	err = row.Scan(&shortLink)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"short": HostURL + shortLink,
			"long":  req.Link,
		})
		return
	}

	// Проверка на наличие короткой ссылки в БД и генерация
	var shortLinkCheck string
	for {
		b := make([]byte, 6)
		rand.Read(b)
		shortLink = base64.URLEncoding.EncodeToString(b)[:6]
		
		row = h.db.QueryRow(c, "SELECT short_link FROM links WHERE short_link = $1", shortLink)
		err = row.Scan(&shortLinkCheck)
		if errors.Is(err, pgx.ErrNoRows) {
			break
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintln("Ошибка в БД: ", err))
			return
		}
	}

	// Добавляем в БД
	_, err = h.db.Exec(c, "INSERT INTO links (long_link, short_link) VALUES ($1, $2)", req.Link, shortLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short": HostURL + shortLink,
		"long":  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")
	var longLink string

	row := h.db.QueryRow(c, "SELECT long_link FROM links WHERE short_link = $1", shortLink)
	err := row.Scan(&longLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, "Ссылка не найдена!")
			return
		}
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже!")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)
}
