package config

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	for _, tst := range []struct {
		expect Config
		in     string
	}{
		{
			expect: Config{
				Handlers: NewHandlerMap(),
			},
			in: "",
		},
		{
			expect: Config{
				Handlers: NewHandlerMap(),
				Sources:  []string{"foo.sh", "bar.sh"},
			},
			in: "include foo.sh\ninclude bar.sh",
		},
		{
			expect: Config{
				Handlers: NewHandlerMap().Add(
					http.MethodGet,
					"/now", "date",
				),
				Sources: []string{"foo.sh"},
			},
			in: "include foo.sh\nget /now date",
		},
		{
			expect: Config{
				Handlers: NewHandlerMap().Add(
					http.MethodGet,
					"/now", "date",
				),
				Sources: []string{"foo.sh"},
			},
			in: "include foo.sh; get /now date",
		},
		{
			expect: Config{
				Handlers: NewHandlerMap().Add(
					http.MethodPost,
					"/docs/:name", "create_doc",
				).Add(
					http.MethodGet,
					"/docs/:name", "read_doc",
				).Add(
					http.MethodPut,
					"/docs/:name", "update_doc",
				).Add(
					http.MethodDelete,
					"/docs/:name", "delete_doc",
				),
			},
			in: `
POST    /docs/:name create_doc
GET     /docs/:name read_doc
PUT     /docs/:name update_doc
DELETE  /docs/:name delete_doc
			`,
		},
	} {
		config, err := Parse(tst.in)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &tst.expect, config)
	}
}
