package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmlbl/echox/shell"
)

func testServer(t *testing.T) *Server {
	server, err := New(1, []string{}, shell.Bash)
	if err != nil {
		t.Error(err)
	}
	return server
}

func testGet(t *testing.T, method, path string) *http.Request {
	r, err := http.NewRequest(method, path, bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Error(err)
	}
	return r
}

func TestNotFound(t *testing.T) {
	server := testServer(t)
	w := httptest.NewRecorder()
	r := testGet(t, http.MethodGet, "/foo")
	server.ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
}
