package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/tmlbl/echox/config"
	"github.com/tmlbl/echox/server"
	"github.com/tmlbl/echox/shell"
)

func main() {
	cfg := config.New()

	// If the arguments are config files, load them
	for _, arg := range os.Args[1:] {
		if _, err := os.Stat(arg); err == nil {
			fmt.Println("Loading config from", arg)
			data, err := ioutil.ReadFile(arg)
			if err != nil {
				log.Fatalln(err)
			}
			c, err := config.Parse(string(data))
			if err != nil {
				log.Fatalln(err)
			}
			cfg.Merge(c)
		}
	}

	// If we got no arguments, attempt to read config from stdin
	if len(os.Args) == 1 {
		data, err := ioutil.ReadAll(os.Stdin)
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

	fmt.Println("[echox] listening on port 7000")
	http.ListenAndServe(":7000", srv)
}
