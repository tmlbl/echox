package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicMatch(t *testing.T) {
	rt := NewRouteTree()
	rt.Add(http.MethodGet, "/foo/bar", "baz")
	item, _ := rt.Match(http.MethodGet, "/foo/bar")
	assert.Equal(t, "baz", item)
}

func TestMethodMatch(t *testing.T) {
	rt := NewRouteTree()
	rt.Add(http.MethodGet, "/foo", "bar")
	rt.Add(http.MethodPost, "/foo", "baz")
	item, _ := rt.Match(http.MethodGet, "/foo")
	assert.Equal(t, "bar", item)
	item, _ = rt.Match(http.MethodPost, "/foo")
	assert.Equal(t, "baz", item)
}

func TestParamMatch(t *testing.T) {
	rt := NewRouteTree()
	rt.Add(http.MethodGet, "/hello/:name", "hello")
	item, params := rt.Match(http.MethodGet, "/hello/tim")
	assert.Equal(t, "hello", item)
	assert.Equal(t, map[string]string{"name": "tim"}, params)
}

func TestDeepParamMatch(t *testing.T) {
	rt := NewRouteTree()
	rt.Add(http.MethodGet, "/parent/:parent/child/:child", "hello")
	item, params := rt.Match(http.MethodGet, "/parent/foo/child/bar")
	assert.Equal(t, "hello", item)
	assert.Equal(t, map[string]string{
		"parent": "foo",
		"child":  "bar",
	}, params)
}

func TestMultiMethod(t *testing.T) {
	rt := NewRouteTree()
	rt.Add(http.MethodGet, "/foo/bar", "baz")
	rt.Add(http.MethodPost, "/foo/bar", "doot")
	item, _ := rt.Match(http.MethodGet, "/foo/bar")
	assert.Equal(t, "baz", item)
	item, _ = rt.Match(http.MethodPost, "/foo/bar")
	assert.Equal(t, "doot", item)
}
