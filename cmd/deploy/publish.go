package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
)

func Publish(t *task.Task) error {
	path := t.Dir()
	cmd := ""

	// Atomic symlink does not supported.
	// Will use simple≤ two steps switch.
	cmd = fmt.Sprintf("cd %s && mv -T release current", path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	return nil
}
