package config

import (
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
				Handlers: map[string]string{},
			},
			in: "",
		},
		{
			expect: Config{
				Handlers: map[string]string{},
				Sources:  []string{"foo.sh", "bar.sh"},
			},
			in: "include foo.sh\ninclude bar.sh",
		},
		{
			expect: Config{
				Handlers: map[string]string{
					"/now": "date",
				},
				Sources: []string{"foo.sh"},
			},
			in: "include foo.sh\nhandle /now date",
		},
		{
			expect: Config{
				Handlers: map[string]string{
					"/now": "date",
				},
				Sources: []string{"foo.sh"},
			},
			in: "include foo.sh; handle /now date",
		},
	} {
		config, err := Parse(tst.in)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &tst.expect, config)
	}
}
