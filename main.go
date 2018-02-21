package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/Sharykhin/golang-todos/handler"
)

func main() {
	fmt.Println("Listen and serve on port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler.Handler()))
}
