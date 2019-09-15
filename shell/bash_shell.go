package shell

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
)

// BashShell contains the state for one bash process
type BashShell struct {
	c    *exec.Cmd
	in   io.WriteCloser
	out  io.ReadCloser
	lock sync.Mutex
	busy bool
}

// Bash starts a bash process for commands to be run in
func Bash() Shell {
	c := exec.Command("bash")
	in, _ := c.StdinPipe()
	out, _ := c.StdoutPipe()
	s := &BashShell{
		c:   c,
		in:  in,
		out: out,
	}
	s.c.Start()
	return s
}

// Include imports a bash source file into the process
func (s *BashShell) Include(file string) error {
	_, err := s.Exec(fmt.Sprintf("source %s", file), nil)
	return err
}

func (s *BashShell) getOutput(cmd string) ([]byte, error) {
	s.in.Write([]byte(cmd))
	s.in.Write([]byte(";printf \"%03d\" $?;echo echox-end\n"))

	var buf []byte
	for {
		var b = make([]byte, 1)
		s.out.Read(b)
		buf = append(buf, b...)
		if endswith(buf, []byte(endDelim)) {
			break
		}
	}

	delim := len(buf) - len(endDelim)
	status := buf[delim-3 : delim]
	code, _ := strconv.Atoi(string(status))
	if code != 0 {
		return nil, fmt.Errorf("Error %d", code)
	}

	return buf[:delim-3], nil
}

// Define defines a variable inside the bash process
func (s *BashShell) define(name, value string) error {
	_, err := s.in.Write([]byte(fmt.Sprintf("%s='%s'\n", name, value)))
	return err
}

// Exec runs the given command in the bash shell
func (s *BashShell) Exec(cmd string, defs map[string]string) ([]byte, error) {
	s.lock.Lock()
	s.busy = true
	defer s.lock.Unlock()

	// Define temporary variables
	for name, value := range defs {
		s.define(name, value)
	}

	out, err := s.getOutput(cmd)
	if err != nil {
		return nil, err
	}

	s.busy = false
	return out, nil
}

// Busy informs the caller whether a process is already running
func (s *BashShell) Busy() bool {
	return s.busy
}
