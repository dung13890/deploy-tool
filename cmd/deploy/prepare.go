package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
	"time"
)

func Prepare(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("if [ ! -d %s ]; then mkdir -p %s; fi", path, path)
	err := t.Run(cmd)
	if err != nil {
		return err
	}

	// Create metadata .dep dir.
	cmd = fmt.Sprintf("cd %s && if [ ! -d .dep ]; then mkdir .dep; fi", path)
	err = t.Run(cmd)
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

	err = symlink(t)
	if err != nil {
		return err
	}

	return nil
}

func symlink(t *task.Task) error {
	path := t.Dir()
	cmd := ""
	cmd = fmt.Sprintf("cd %s && if [ -h release ]; then rm release; rm -rf $(readlink release); fi", path)
	err := t.Run(cmd)
	if err != nil {
		return err
	}
	now := time.Now()
	folder := fmt.Sprintf("releases/%s", now.Format("20060102150405"))

	// Make new release
	cmd = fmt.Sprintf("cd %s && mkdir -p %s", path, folder)
	err = t.Run(cmd)
	if err != nil {
		return err
	}

	// Make symlink release
	cmd = fmt.Sprintf("cd %s && ln -nfs %s release", path, folder)
	err = t.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}
