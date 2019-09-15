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
	} {
		config, err := Parse(tst.in)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &tst.expect, config)
	}
}
