package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewPMServer(NewInMemoryPMStore())
	log.Fatal(http.ListenAndServe(":4673", server))
}
