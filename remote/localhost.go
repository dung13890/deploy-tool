package remote

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Localhost struct {
	cmd     *exec.Cmd
	user    string
	group   string
	dir     string
	project string
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
	running bool
}

func (l *Localhost) Load(_ string, _ string, group string, _ int, dir string, project string) {
	u, _ := user.Current()
	l.user = u.Username
	l.dir = dir
	l.project = project
	l.group = group
}

func (l *Localhost) GetDirectory() string {
	return filepath.Join(l.dir, l.project)
}

func (l *Localhost) GetUser() (string, string) {
	return l.user, l.group
}

func (l *Localhost) Prefix() string {
	return fmt.Sprintf("[%s@localhost]", l.user)
}

func (l *Localhost) Connect(_ string) error {
	return nil
}

func (l *Localhost) Run(cmd string) error {
	if l.running {
		return errors.New("Command already running")
	}
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

	if err := l.cmd.Start(); err != nil {
		return err
	}

	l.running = true

	return nil
}

func (l *Localhost) Wait() error {
	if !l.running {
		return errors.New("Trying to wait on stopped command")
	}
	err := l.cmd.Wait()
	l.running = false

	return err
}

func (l *Localhost) CombinedOutput(cmd string) (out []byte, err error) {
	l.cmd = exec.Command("bash", "-c", cmd)

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
