package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
	"strings"
	"time"
)

func Prepare(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("if [ ! -d %s ]; then mkdir -p %s; fi", path, path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	// Create metadata .dep dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d .dep ]; then mkdir .dep; fi", path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	// Create releases dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d releases ]; then mkdir releases; fi", path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	// Create shared dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d shared ]; then mkdir shared; fi", path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	if err := symlink(t); err != nil {
		return err
	}

	if err := cleanup(t); err != nil {
		return err
	}

	return nil
}

func symlink(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("cd %s && if [ -h release ]; then rm release; rm -rf $(readlink release); fi", path)
	if err := t.Run(cmd); err != nil {
		return err
	}
	now := time.Now()
	folder := fmt.Sprintf("releases/%s", now.Format("20060102150405"))

	// Make new release
	cmd = fmt.Sprintf("cd %s && mkdir -p %s", path, folder)
	if err := t.Run(cmd); err != nil {
		return err
	}

	// Make symlink release
	cmd = fmt.Sprintf("cd %s && ln -nfs %s release", path, folder)
	if err := t.Run(cmd); err != nil {
		return err
	}

	return nil
}

func cleanup(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("cd %s/releases && ls -t -1 -d */", path)
	// Will list only dirs in releases.
	out, err := t.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	out = strings.Replace(strings.TrimSpace(out), "\r\n", "\n", -1)
	arr := strings.Split(out, "\n")
	if len(arr) > 5 {
		for _, v := range arr[5:] {
			v = strings.TrimSuffix(v, "/")
			cmd = fmt.Sprintf("cd %s/releases && rm -rf %s", path, v)
			if err := t.Run(cmd); err != nil {
				return err
			}
		}
	}

	return err
}
