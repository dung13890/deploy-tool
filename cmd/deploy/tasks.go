package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/utils"
	"strings"
)

type Tasks struct {
	list []string
}

func NewTasks(list []string) *Tasks {
	return &Tasks{
		list: list,
	}
}

func (ts *Tasks) Run(t *task.Task) error {
	path := t.GetDirectory()
	releasePath := fmt.Sprintf("%s/release", path)
	cmd := ""

	// Loop list unique tasks
	for _, v := range utils.UniqueArr(ts.list) {
		v = strings.TrimSpace(v)
		cmd = fmt.Sprintf("cd %s && %s", releasePath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}
