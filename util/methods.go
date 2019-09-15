package util

import "net/http"

// HTTPMethods is a list of every method
var HTTPMethods = []string{
	http.MethodGet,
	http.MethodPost, http.MethodConnect, http.MethodDelete, http.MethodHead,
	http.MethodOptions, http.MethodPatch, http.MethodPut, http.MethodTrace,
}
