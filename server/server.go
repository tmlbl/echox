package server

import (
	"fmt"
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

	out, err := process.Exec(item, params)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write(out)
	}
}
