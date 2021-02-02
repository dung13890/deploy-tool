package remote

import (
	"io"
)

type Remote interface {
	Load(address string, user string, port int, dir string)
	Dir() string
	Connect(pathKey string) error
	Run(cmd string) error
	Stdout() io.Reader
	StdErr() io.Reader
	Close() error
}
