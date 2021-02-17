package remote

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

type Localhost struct {
	cmd    *exec.Cmd
	user   string
	dir    string
	stdin  io.WriteCloser
	stdout io.Reader
	stderr io.Reader
}

func (l *Localhost) Load(_ string, _ string, _ int, dir string) {
	u, _ := user.Current()
	l.user = u.Username
	l.dir = dir
}

func (l *Localhost) Dir() string {
	return fmt.Sprintf("%s/data/sites/%s", os.Getenv("HOME"), l.dir)
}

func (l *Localhost) Prefix() string {
	return fmt.Sprintf("[%s@localhost]", l.user)
}

func (l *Localhost) Connect(_ string) error {
	return nil
}

func (l *Localhost) Run(cmd string) error {
	l.cmd = exec.Command("bash", "-c", cmd)
	var err error

	if l.stdin, err = l.cmd.StdinPipe(); err != nil {
		return err
	}

	if l.stdout, err = l.cmd.StdoutPipe(); err != nil {
		return err
	}

	if l.stderr, err = l.cmd.StderrPipe(); err != nil {
		return err
	}

	if err := l.cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (l *Localhost) CombinedOutput(cmd string) (out []byte, err error) {
	l.cmd = exec.Command("bash", "-c", cmd)

	if l.stderr, err = l.cmd.StderrPipe(); err != nil {
		return
	}

	out, err = l.cmd.CombinedOutput()
	return
}

func (l *Localhost) Stdin() io.WriteCloser {
	return l.stdin
}

func (l *Localhost) Stdout() io.Reader {
	return l.stdout
}

func (l *Localhost) StdErr() io.Reader {
	return l.stderr
}

func (l *Localhost) Close() error {
	return nil
}
