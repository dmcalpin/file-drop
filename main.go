package main

import (
	"fmt"
	"log"
	"net/http"

	server "github.com/dmcalpin/file-drop/server"
)

func main() {
	port := ":9988"
	fmt.Println("Starting server:", "http://localhost"+port)

	mux := server.NewDefaultMux()
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
