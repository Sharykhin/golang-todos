package main

import (
	"log"
	"net/http"

	"github.com/Sharykhin/golang-todos/handler"
)

// @QUESTION:
// Do we need in some sort of config package?

func main() {
	log.Fatal(http.ListenAndServe(":8082", handler.Handler()))
}
