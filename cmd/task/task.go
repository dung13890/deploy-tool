package task

import (
	"github.com/dung13890/deploy-tool/remote"
	"io"
	"log"
	"os"
	"sync"
)

type Task struct {
	remote remote.Remote
	debug  bool
	cmd    string
}

func NewTask(r remote.Remote, d bool) *Task {
	return &Task{
		remote: r,
		debug:  d,
	}
}

func (t *Task) printLog() error {
	wg := sync.WaitGroup{}
	// Copy over tasks's STDOUT.
	wg.Add(1)
	go func(r remote.Remote) {
		defer wg.Done()
		_, err := io.Copy(os.Stdout, r.Stdout())
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}(t.remote)
	// Copy over tasks's STDERR.
	wg.Add(1)
	go func(r remote.Remote) {
		defer wg.Done()
		_, err := io.Copy(os.Stderr, r.StdErr())
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}(t.remote)
	wg.Wait()

	return nil
}

func (t *Task) GetDirectory() string {
	return t.remote.GetDirectory()
}

func (t *Task) GetUser() (string, string) {
	return t.remote.GetUser()
}

func (t *Task) Run(cmd string) error {
	t.cmd = cmd
	if t.debug {
		t.cmd = "set -x;" + cmd
	}
	if err := t.remote.Run(t.cmd); err != nil {
		return err
	}
	if t.debug {
		t.printLog()
	}

	if err := t.remote.Wait(); err != nil {
		return err
	}

	return nil
}

func (t *Task) CombinedOutput(cmd string) (out string, err error) {
	t.cmd = cmd
	o, err := t.remote.CombinedOutput(t.cmd)
	out = string(o)

	return
}
