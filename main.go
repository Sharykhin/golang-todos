package main

import (
	"log"
	"net/http"

	"github.com/Sharykhin/golang-todos/handler"
)

func main() {
	log.Fatal(http.ListenAndServe(":8082", handler.Handler()))
}
