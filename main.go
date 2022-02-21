package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	fileName = "pm.storage.json"
	port     = 4673
)

func main() {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("Couldn't open the '%s' '%v'", fileName, err)
	}

	store, err := NewFSStore(file)

	if err != nil {
		log.Fatalf("Couldn't create file system store '%v'", err)
	}
	srv := NewPMServer(store)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), srv); err != nil {
		log.Fatalf("Couldn't listen on port %d %v", port, err)
	}
}
