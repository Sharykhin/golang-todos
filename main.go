package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Sharykhin/golang-todos/command"
	"github.com/Sharykhin/golang-todos/handler"
)

func main() {
	migrate := flag.Bool("migrate", false, "whether to run migrations or not")
	flag.Parse()
	if *migrate == true {
		command.Migrate()
	}
	fmt.Println("Listen and serve on port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler.Handler()))
}
