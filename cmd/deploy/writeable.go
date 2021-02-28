package deploy

import (
	"fmt"
	"github.com/dung13890/deploy-tool/cmd/task"
	"strings"
)

func Writeable(t *task.Task) (err error) {
	path := t.GetDirectory()
	user, group := t.GetUser()
	if group == "" {
		group, err = getGroupUser(t)
		if err != nil {
			return
		}
	}

	cmd := fmt.Sprintf("sudo chown -L -R %s:%s %s", user, group, path)
	if err := t.Run(cmd); err != nil {
		return err
	}

	return
}

func getGroupUser(t *task.Task) (string, error) {
	groupUser := "$( id -gn )"
	cmd := "ps axo comm,user | grep -E '[a]pache|[h]ttpd|[_]www|[w]ww-data|[n]ginx' | grep -v root | sort | awk '{print $NF}' | uniq"
	out, err := t.CombinedOutput(cmd)
	if err != nil {
		return groupUser, err
	}
	out = strings.Replace(strings.TrimSpace(out), "\r\n", "\n", -1)
	arr := strings.Split(out, "\n")
	if len(arr) > 0 && arr[0] != "" {
		return arr[0], nil
	}

	return groupUser, nil
}
