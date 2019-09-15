package config

import (
	"net/http"
	"strings"
)

var httpMethods = []string{
	http.MethodGet, http.MethodConnect, http.MethodDelete, http.MethodHead,
	http.MethodOptions, http.MethodPatch, http.MethodPut, http.MethodTrace,
}

func isMethod(s string) bool {
	for _, m := range httpMethods {
		if s == m {
			return true
		}
	}
	return false
}

type HandlerMap map[string]map[string]string

func NewHandlerMap() HandlerMap {
	h := HandlerMap{}
	for _, m := range httpMethods {
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
	for k, v := range c2.Handlers {
		c.Handlers[k] = v
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
		args := strings.Split(line, " ")
		cmd := strings.ToUpper(args[0])
		if cmd == "INCLUDE" {
			cfg.Sources = append(cfg.Sources, args[1])
		} else if isMethod(cmd) {
			cfg.Handlers.Add(cmd, args[1], args[2])
		}
	}

	return &cfg, nil
}
