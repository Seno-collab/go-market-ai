package main

import (
	_ "go-ai/docs"
	"log"

	"go-ai/internal/app"
)

// go ai
// mission using golang build ai. Development application AI

func main() {
	e := app.NewServer()
	if err := app.Run(e); err != nil {
		log.Fatal(err)
	}
}
