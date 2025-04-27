package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleGet(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello, %s! This is a GET request.", name)

}

func handlePost(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	fmt.Fprintf(w, "Received POST data: %s", string(body))
}

// 在一个 goroutine 中启动 HTTP 服务器
// go startServer()

func startServer() {

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/post", handlePost)

	if err := http.ListenAndServe(":"+string(rune(_appConfig.GuiPort)), nil); err != nil {
		log.Fatal(err)
	}
}
