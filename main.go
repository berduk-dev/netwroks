package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Ошибка. Статус код = %d\n", resp.StatusCode)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Не смог прочитать: %v\n", err)
		return
	}

	var posts []Post

	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		fmt.Println("Не удалось анмаршалить: ", err)
	}

	fmt.Printf("%+v", posts)
}
