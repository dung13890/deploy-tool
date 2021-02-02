package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/remote"
	"log"
)

func Prepare(remote remote.Remote) error {
	t := &task{r: remote}
	err := t.setup()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (t *task) setup() error {
	path := t.r.Dir()
	err := t.r.Run(fmt.Sprintf("ls -al && if [ ! -d %s ]; then mkdir -p %s; fi ", path, path))
	if err != nil {
		return err
	}
	// Create releases dir.
	err = t.r.Run(fmt.Sprintf("ls -al && cd %s && if [ ! -d releases ]; then mkdir releases; fi", path))
	if err != nil {
		return err
	}
	// Create shared dir.
	err = t.r.Run(fmt.Sprintf("ls -al && cd %s && if [ ! -d shared ]; then mkdir shared; fi", path))
	if err != nil {
		return err
	}

	return nil
}
