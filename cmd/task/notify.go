package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Notify struct {
	address         string
	project         string
	token           string
	room            string
	to              string
	slackWebhook    string
	otherUrlWebhook string
	otherChannel    string
	feature         string
}

func NewNotify(
	address string,
	project string,
	token string,
	room string,
	to string,
	slackWebhook string,
	feature string,
	otherUrlWebhook string,
	otherChannel string,
) *Notify {
	return &Notify{
		address:         address,
		project:         project,
		token:           token,
		room:            room,
		to:              to,
		slackWebhook:    slackWebhook,
		feature:         feature,
		otherUrlWebhook: otherUrlWebhook,
		otherChannel:    otherChannel,
	}
}

func (n *Notify) Push(status string) error {
	if n.token != "" {
		n.doSendChatwork(status)
	}
	if n.slackWebhook != "" {
		n.doSendSlack(status)
	}
	if n.otherUrlWebhook != "" {
		n.doSendOther(status)
	}

	return nil
}

func (n *Notify) doSendSlack(status string) (resq []byte, err error) {

	// Make message
	msg := []byte(fmt.Sprintf("{'text': '*Deploy (%s) into Server (%s)*```Build Status: %s\n%s```'}",
		n.project,
		n.address,
		status,
		n.feature,
	))

	// Make request
	client := &http.Client{}

	req, err := http.NewRequest("POST", n.slackWebhook, bytes.NewBuffer(msg))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Execute request
	res, err := client.Do(req)
	defer res.Body.Close()

	resq, _ = ioutil.ReadAll(res.Body)

	return
}

func (n *Notify) doSendOther(status string) (resq []byte, err error) {
	// Make message
	msg := fmt.Sprintf("{'text': '*Deploy (%s) into Server (%s)*```Build Status: %s\n%s```'}",
		n.project,
		n.address,
		status,
		n.feature,
	)

	postBody := []byte(fmt.Sprintf("{'service': 'slack', 'channel': %s, 'receivers': 'here', 'message': %s}",
		n.otherChannel,
		msg,
	))

	// Make request
	client := &http.Client{}

	req, err := http.NewRequest("POST", n.otherUrlWebhook, bytes.NewBuffer(msg))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Execute request
	res, err := client.Do(req)
	defer res.Body.Close()

	resq, _ = ioutil.ReadAll(res.Body)

	return
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
