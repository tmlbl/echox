package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/tmlbl/echox/config"
	"github.com/tmlbl/echox/server"
	"github.com/tmlbl/echox/shell"
)

func main() {
	cfg := config.New()

	// If the arguments are config files, load them
	for _, arg := range os.Args[1:] {
		if _, err := os.Stat(arg); err == nil {
			data, err := os.ReadFile(arg)
			if err != nil {
				log.Fatalln(err)
			}
			if strings.HasSuffix(arg, ".sh") {
				c, err := config.ParseBash(string(data))
				if err != nil {
					log.Fatalln(err)
				}
				c.Sources = append(c.Sources, arg)
				cfg.Merge(c)
			} else {
				c, err := config.Parse(string(data))
				if err != nil {
					log.Fatalln(err)
				}
				cfg.Merge(c)
			}
		}
	}

	// If we got no arguments, attempt to read config from stdin
	if len(os.Args) == 1 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		c, err := config.Parse(string(data))
		cfg.Merge(c)
	}

	// Show include files / make sure they exist
	for _, src := range cfg.Sources {
		if _, err := os.Stat(src); err != nil {
			fmt.Printf("ERROR reading %s: %s\n", src, err)
		}
	}

	// Instantiate the server
	srv, err := server.New(runtime.NumCPU(), cfg.Sources, shell.Bash)
	if err != nil {
		log.Fatalln(err)
	}

	// Set up handlers for the config
	for method := range cfg.Handlers {
		for path, cmd := range cfg.Handlers[method] {
			fmt.Println("[echox]", method, path, "=>", cmd)
			srv.AddRoute(method, path, cmd)
		}
	}

	addr := ":7171"

	fmt.Printf("[echox] listening on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, srv))
}
