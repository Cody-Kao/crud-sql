package main

import (
	"log"

	"github.com/Cody-Kao/crud-sql/server"
)

func main() {
	server := server.CreateServer()
	log.Fatal(server.ListenAndServe())
}
