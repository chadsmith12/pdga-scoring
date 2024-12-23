package main

import (
	"log"

	"github.com/chadsmith12/pdga-scoring/internal/server"
)

func main() {
	server := server.NewServer()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
