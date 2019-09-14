package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tmlbl/echox/config"
	"github.com/tmlbl/echox/shell"
)

func handler(sources []string, cmd string) func(w http.ResponseWriter, r *http.Request) {
	ps, err := shell.Bash()
	if err != nil {
		panic(err)
	}
	for _, src := range sources {
		ps.Include(src)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Supply headers to the function in a format amenable to shell
		// parsing
		headers := []string{}
		for k, v := range r.Header {
			headers = append(headers,
				fmt.Sprintf("%s:%s", k, strings.Join(v, ",")))
		}
		hdrs := strings.Join(headers, "\n")

		out, err := ps.Exec(cmd, map[string]string{
			"headers": hdrs,
		})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(200)
			w.Write(out)
		}
	}
}

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

	mux := http.NewServeMux()

	// Set up handlers for the config
	for path, cmd := range cfg.Handlers {
		fmt.Println("[echox]", path, "=>", cmd)
		mux.HandleFunc(path, handler(cfg.Sources, cmd))
	}

	fmt.Println("[echox] listening on port 7000")
	http.ListenAndServe(":7000", mux)
}
