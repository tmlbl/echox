package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tmlbl/echox/shell"
)

// Server handles incoming HTTP requests and forwards them to the appropriate
// handler
type Server struct {
	routes *RouteTree
	shells []shell.Shell
}

// New instantiates a new Server
func New(nproc int, sources []string, provider func() shell.Shell) (*Server, error) {
	s := Server{
		routes: NewRouteTree(),
	}
	fmt.Printf("[echox] starting %d workers\n", nproc)
	for i := 0; i < nproc; i++ {
		sh := provider()
		for _, src := range sources {
			sh.Include(src)
		}
		s.shells = append(s.shells, sh)
	}
	return &s, nil
}

// AddRoute adds an entry to the RouteTree
func (s *Server) AddRoute(method, path, item string) {
	s.routes.Add(method, path, item)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	item, params := s.routes.Match(r.Method, r.URL.Path)
	if item == "" {
		// No handler found
		w.WriteHeader(404)
		return
	}

	// Find an available process
	var process shell.Shell
	for _, sh := range s.shells {
		if !sh.Busy() {
			process = sh
			break
		}
	}
	if process == nil {
		// Service not available?
		w.WriteHeader(503)
		w.Write([]byte("All processes are busy"))
		return
	}

	// Supply headers to the function in a format amenable to shell
	// parsing
	headers := []string{}
	for k, v := range r.Header {
		headers = append(headers,
			fmt.Sprintf("%s:%s", k, strings.Join(v, ",")))
	}
	params["headers"] = strings.Join(headers, "\n")

	// Supply body content in the body variable
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err == nil {
			params["body"] = string(body)
		}
	}

	fmt.Printf("[echox] %s %s\n", r.Method, r.URL.Path)
	out, err := process.Exec(item, params)

	// Write a CORS header because CORS is annoying
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write(out)
	}
}
