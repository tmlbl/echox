package config

import (
	"strings"
)

// Config represents config information required to run a server
type Config struct {
	Sources  []string
	Handlers map[string]string
}

// New creates a new Config instance with a non-nil Handlers map
func New() *Config {
	return &Config{
		Handlers: map[string]string{},
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
		Handlers: map[string]string{},
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
		switch args[0] {
		case "include":
			cfg.Sources = append(cfg.Sources, args[1])
		case "handle":
			cfg.Handlers[args[1]] = args[2]
		}
	}

	return &cfg, nil
}
