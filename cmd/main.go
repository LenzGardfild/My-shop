package main

import (
	"log"
	"net/http"

	"my-shop/internal/server"
)

func main() {
	srv := server.New()
	log.Println("HTTP сервер запущен на :8081")
	http.ListenAndServe(":8081", srv)
}
