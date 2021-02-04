package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
)

func Prepare(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("if [ ! -d %s ]; then mkdir -p %s; fi", path, path)
	err := t.Run(cmd)
	if err != nil {
		return err
	}

	// Create releases dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d releases ]; then mkdir releases; fi", path)
	err = t.Run(cmd)
	if err != nil {
		return err
	}

	// Create shared dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d shared ]; then mkdir shared; fi", path)
	err = t.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}
