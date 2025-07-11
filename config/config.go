package config

import (
	"strings"

	"github.com/tmlbl/echox/util"
)

func isMethod(s string) bool {
	for _, m := range util.HTTPMethods {
		if s == m {
			return true
		}
	}
	return false
}

type HandlerMap map[string]map[string]string

func NewHandlerMap() HandlerMap {
	h := HandlerMap{}
	for _, m := range util.HTTPMethods {
		h[m] = map[string]string{}
	}
	return h
}

func (h HandlerMap) Add(method, key, value string) HandlerMap {
	h[method][key] = value
	return h
}

// Config represents config information required to run a server
type Config struct {
	Sources  []string
	Handlers HandlerMap
}

// New creates a new Config instance with a non-nil Handlers map
func New() *Config {
	return &Config{
		Handlers: NewHandlerMap(),
	}
}

// Merge combines another Config into the receiver
func (c *Config) Merge(c2 *Config) {
	c.Sources = append(c.Sources, c2.Sources...)
	for method, handlers := range c2.Handlers {
		for k, v := range handlers {
			c.Handlers[method][k] = v
		}
	}
}

// Parse parses the input string into a config object
func Parse(s string) (*Config, error) {
	cfg := Config{
		Handlers: NewHandlerMap(),
	}

	lines := strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == ';'
	})
	for _, line := range lines {
		line = strings.Trim(line, " \t")
		if len(line) == 0 {
			continue
		}
		args := strings.Fields(line)
		cmd := strings.ToUpper(args[0])
		if cmd == "INCLUDE" {
			cfg.Sources = append(cfg.Sources, args[1])
		} else if isMethod(cmd) {
			cfg.Handlers.Add(cmd, args[1], args[2])
		}
	}

	return &cfg, nil
}

// ParseBash parses a bash script and builds a handler map
// by looking for comments in this format:
// # GET /greet/:name
// Above function definitions
func ParseBash(s string) (*Config, error) {
	cfg := Config{
		Handlers: NewHandlerMap(),
	}

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 || parts[0] != "#" {
			continue
		}

		method := strings.ToUpper(parts[1])
		if !isMethod(method) {
			continue
		}
		path := parts[2]

		if i+1 >= len(lines) {
			continue
		}
		nextLine := strings.TrimSpace(lines[i+1])

		// must be a function definition e.g. `my_func()`
		// and not a variable assignment e.g. `my_var=(...)`
		if !strings.Contains(nextLine, "()") {
			continue
		}
		parenIx := strings.Index(nextLine, "(")
		if parenIx <= 0 {
			continue
		}

		cmd := strings.TrimSpace(nextLine[:parenIx])
		cmd = strings.TrimPrefix(cmd, "function")
		cmd = strings.TrimSpace(cmd)

		if cmd == "" {
			continue
		}

		cfg.Handlers.Add(method, path, cmd)
	}
	return &cfg, nil
}
