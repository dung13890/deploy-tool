package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	passVersion = "1.2"
	endpoint    = "http://127.0.0.1:1308/openssl"
	key         = "r3kqouWSZ1"
)

type openssl struct {
	config     config.Configuration
	privateKey string
	log        bool
}

type paramRequest struct {
	Project string `json:"project"`
	Value   string `json:"value"`
	Result  bool   `json:"result"`
}

func NewOpenssl() *cli.Command {
	return &cli.Command{
		Name:    "openssl",
		Aliases: []string{"o"},
		Usage:   "Check the SSL/TLS Cipher Suites in Servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
			config.EnableLog,
		},
		Action: func(ctx *cli.Context) error {
			o := &openssl{}
			err := o.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			o.log = ctx.Bool("log")
			o.privateKey = ctx.String("identity")
			o.exec()
			return nil
		},
	}
}

func (o *openssl) exec() error {
	var r remote.Remote

	if o.config.Server.Address == "127.0.0.1" || o.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}
	defer r.Close()
	r.Load(
		o.config.Server.Address,
		o.config.Server.User,
		o.config.Server.Group,
		o.config.Server.Port,
		o.config.Server.Dir,
		o.config.Server.Project,
	)
	fmt.Println("Check the SSL/TLS Cipher Suites in Servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" %s: Processing...", r.Prefix())
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s %s: Completed!\n", green("✔"), r.Prefix())
	sp.Start()
	if err := r.Connect(o.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
	}
	t := task.NewTask(r, o.log)
	rs, err := o.checkOpenssl(t)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	sp.Stop()
	if rs == true {
		success := color.New(color.FgHiGreen, color.Bold).PrintlnFunc()
		success("[PASS] Server used TLS with version >= 1.2!")
		return nil
	}
	fail := color.New(color.FgHiRed, color.Bold).PrintlnFunc()
	fail("[FAILED] Please install TLS with version >= 1.2")

	return nil
}

func (o *openssl) checkOpenssl(t *task.Task) (pass bool, err error) {
	pass = false
	cmd := "openssl ciphers -v | awk '{print $2}' | sort | uniq"
	// Run Command check openssl
	out, err := t.CombinedOutput(cmd)
	if err != nil {
		return
	}
	// Handler output
	out = strings.Replace(strings.TrimSpace(out), "\r\n", "\n", -1)
	arr := strings.Split(out, "\n")
	for _, v := range arr {
		regEx := regexp.MustCompile(`TLSv(\d+(?:\.\d+)+)`)
		rs := regEx.FindStringSubmatch(v)
		if len(rs) > 1 {
			f, err := strconv.ParseFloat(rs[1], 8)
			if err != nil {
				return false, err
			}
			if f >= 1.2 {
				pass = true
			}
		}
	}

	params := &paramRequest{
		Project: o.config.Server.Project,
		Value:   strings.Join(arr[:], "-"),
		Result:  pass,
	}
	o.notify(params)
	return
}

func (o *openssl) notify(params *paramRequest) error {
	// Make request
	client := &http.Client{}
	p, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(p))
	req.Header.Add("X-SecretKey", key)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	if _, err = client.Do(req); err != nil {
		return err
	}
	return nil
}
