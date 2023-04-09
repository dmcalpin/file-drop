package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// HandleStatus should be called by the client to verify
// the server is up and running.
// This is useful if a client wants to scan the network
// for multiple clients
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET allowed."))
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Fprintf(w, `{"ok": true, "name": "%s"}`, hostname)
}
