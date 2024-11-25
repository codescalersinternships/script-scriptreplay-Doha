package main

import (
	"log"
	"os"

	"github.com/codescalersinternships/script-scriptreplay-Doha/internal/script"
)

var (
	filename = "typescript"
)

func main() {
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	if err := script.Script(filename); err != nil {
		log.Fatal(err)
	}
}
