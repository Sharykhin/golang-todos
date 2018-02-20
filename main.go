package main

import (
	"log"
	"net/http"

	"github.com/Sharykhin/golang-todos/handler"
	"fmt"
)

func main() {
	fmt.Println("Listen and serve on port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler.Handler()))
}
