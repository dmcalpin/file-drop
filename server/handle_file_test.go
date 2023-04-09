package server

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFile(t *testing.T) {
	mux := NewDefaultMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/file", strings.NewReader("test file"))

	mux.ServeHTTP(rr, req)

	if rr.Body.String() != "" {
		t.Error(rr.Body.String(), "fail")
	}

	// check rr
}
