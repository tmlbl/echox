package server

import (
	"strings"

	"github.com/tmlbl/echox/util"
)

// RouteTree stores the association between route parts and handlers
type RouteTree struct {
	part     string
	children []*RouteTree
	item     string
	method   string
}

func NewRouteTree() *RouteTree {
	children := []*RouteTree{}
	for _, method := range util.HTTPMethods {
		children = append(children, &RouteTree{
			method: method,
		})
	}
	return &RouteTree{
		children: children,
	}
}

func (t *RouteTree) getChild(method, part string) *RouteTree {
	for _, c := range t.children {
		if c.part == part && c.method == method {
			return c
		} else if c.method == method && isWildCard(c.part) {
			return c
		}
	}
	return nil
}

func isWildCard(part string) bool {
	return len(part) > 0 && part[0] == ':'
}

func (t *RouteTree) Add(method, path, item string) {
	parts := strings.Split(path, "/")
	tree := t
	for _, part := range parts {
		if part == "" {
			continue
		}
		c := tree.getChild(method, part)
		if c == nil {
			tree.children = append(tree.children, &RouteTree{
				part:   part,
				method: method,
			})
		}
		tree = tree.getChild(method, part)
	}
	tree.item = item
	tree.method = method
}

func (t *RouteTree) Match(method, path string) (string, map[string]string) {
	parts := strings.Split(path, "/")
	tree := t
	params := map[string]string{}

	for _, part := range parts {
		if part == "" {
			continue
		}
		tree = tree.getChild(method, part)
		if tree == nil {
			return "", nil
		}

		if isWildCard(tree.part) {
			params[tree.part[1:]] = part
		}
	}
	return tree.item, params
}
