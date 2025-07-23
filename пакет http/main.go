package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	post1 := Post{
		UserID: 1,
		ID:     321,
		Title:  "Первый пост",
		Body:   "The first post",
	}

	marshalledPost, err := json.Marshal(post1)
	if err != nil {
		log.Printf("Ошибка при маршале: %w", err)
	}

	url := "https://httpbin.org/post"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalledPost))
	if err != nil {
		log.Printf("Ошибка при пост запросе: %w", err)
	}

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении респонса: %w", err)
	}

	fmt.Println(string(bytesResp))
}
