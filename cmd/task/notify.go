package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Notify struct {
	address string
	project string
	token   string
	room    string
	to      string
	feature string
}

func NewNotify(address string, project string, token string, room string, to string, feature string) *Notify {
	return &Notify{
		address: address,
		project: project,
		token:   token,
		room:    room,
		to:      to,
		feature: feature,
	}
}

func (n *Notify) Push(status string) error {
	if n.token != "" {
		n.doSendChatwork(status)
	}

	return nil
}

func (n *Notify) doSendChatwork(status string) (resq []byte, err error) {
	var body string

	to := n.to
	if to == "" {
		to = "[toall]"
	}

	// Make message
	body = fmt.Sprintf(
		"%s\n[info][title]Deploy (%s) into Server (%s)[/title]Build Status: %s\n%s[/info]",
		to,
		n.project,
		n.address,
		status,
		n.feature,
	)
	// Make request
	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%s/messages?body=%s", n.room, url.QueryEscape(body))
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return
	}
	req.Header.Add("X-ChatWorkToken", n.token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	res, err := client.Do(req)
	defer res.Body.Close()

	resq, _ = ioutil.ReadAll(res.Body)

	return
}
