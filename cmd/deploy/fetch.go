package deploy

import (
	"errors"
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
)

type Repo struct {
	url    string
	branch string
	tag    string
}

func NewRepo(url string, branch string, tag string) *Repo {
	return &Repo{
		url:    url,
		branch: branch,
		tag:    tag,
	}
}

func (r *Repo) Fetch(t *task.Task) (err error) {
	path := t.GetDirectory()
	releasePath := fmt.Sprintf("%s/release", path)
	// Update code at release_path on host.
	cmd, err := r.makeCmd(releasePath)
	if err != nil {
		return
	}
	return t.Run(cmd)
}

func (r *Repo) makeCmd(path string) (cmd string, err error) {
	if r.url == "" {
		err = errors.New("empty url of the repository! Please preview your config.yml and try again later")
		return
	}
	at := ""
	// If option `branch` is set
	if r.branch != "" {
		at = fmt.Sprintf("-b %s", r.branch)
	}

	// If option `tag` is set
	if r.tag != "" {
		at = fmt.Sprintf("-b %s", r.tag)
	}

	cmd = fmt.Sprintf("git clone %s %s %s 2>&1", at, r.url, path)

	return
}
