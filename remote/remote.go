package remote

import (
	"io"
)

type Remote interface {
	Load(address string, user string, port int, dir string)
	Dir() string
	Prefix() string
	Connect(pathKey string) error
	Run(cmd string) error
	CombinedOutput(cmd string) (out []byte, err error)
	Stdin() io.WriteCloser
	Stdout() io.Reader
	StdErr() io.Reader
	Close() error
}
