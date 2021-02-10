package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/utils"
	"strings"
)

type Shared struct {
	folders []string
	files   []string
}

func NewShared(folders []string, files []string) *Shared {
	return &Shared{
		folders: folders,
		files:   files,
	}
}

func (s *Shared) Setup(t *task.Task) error {
	path := t.Dir()
	sharedPath := fmt.Sprintf("%s/shared", path)
	releasePath := fmt.Sprintf("%s/release", path)

	//Shared Folders
	for _, v := range utils.UniqueArr(s.folders) {
		v = strings.TrimSpace(strings.TrimPrefix(v, "/"))
		// Check if shared folder does not exist then make shared folder
		cmd := fmt.Sprintf("if [ ! -d %s/%s ]; then mkdir -p %s/%s; fi", sharedPath, v, sharedPath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}

		// Remove from source.
		cmd = fmt.Sprintf("rm -rf %s/%s", releasePath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}

		// Symlink shared dir to release dir
		cmd = fmt.Sprintf("ln -nfs %s/%s %s/%s", sharedPath, v, releasePath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}
	}

	//Shared Files
	for _, v := range utils.UniqueArr(s.files) {
		v = strings.TrimSpace(strings.TrimPrefix(v, "/"))
		// Create dirname of shared file if not existing
		cmd := fmt.Sprintf("if [ ! -d $( dirname %s/files/%s ) ]; then mkdir -p $( dirname %s/files/%s ); fi", sharedPath, v, sharedPath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}

		// Check if shared folder does not exist then make shared folder
		cmd = fmt.Sprintf("if [ ! -f %s/files/%s ]; then touch %s/files/%s; fi", sharedPath, v, sharedPath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}
		// Copy share folder into release.
		cmd = fmt.Sprintf("cp %s/files/%s %s/%s", sharedPath, v, releasePath, v)
		if err := t.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}
