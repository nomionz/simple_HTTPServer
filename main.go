package main

import (
	"log"
	"net/http"
)

func main() {
	server := &ProjectManagementServer{NewInMemoryPMStore()}
	log.Fatal(http.ListenAndServe(":4673", server))
}
