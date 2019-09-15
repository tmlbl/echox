package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashShell(t *testing.T) {
	bash := Bash()
	out, err := bash.Exec("echo hello", nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "hello\n", string(out))
}

func TestBashShellInclude(t *testing.T) {
	bash := Bash()
	err := bash.Include("test/my_func.sh")
	if err != nil {
		t.Error(err)
	}
	out, err := bash.Exec("my_func", nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "my function\n", string(out))
}

func TestBashShellDefs(t *testing.T) {
	bash := Bash()
	out, err := bash.Exec("echo $foo", map[string]string{
		"foo": "bar",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "bar\n", string(out))
}

func TestBashShellError(t *testing.T) {
	bash := Bash()
	out, err := bash.Exec("forkjgorkgg", nil)
	if err == nil {
		t.Errorf("Should have produced an error")
	}
	assert.Equal(t, len(out), 0)
}
