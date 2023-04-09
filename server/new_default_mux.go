package server

import (
	"net/http"
)

// NewDefaultMux sets up the default server mux
// for this server. This is useful when
// running this code via the included main.go file.
// Alternatively, if including this code as an
// import, you can skip this and configure
// your own mux.
func NewDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", HandleStatus)
	mux.HandleFunc("/file", HandleFile)

	return mux
}
