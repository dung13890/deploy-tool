package remote

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	address    string
	user       string
	port       int
	dir        string
	conn       *ssh.Client
	connOpened bool
	stdin      io.WriteCloser
	stdout     io.Reader
	stderr     io.Reader
}

func (s *Server) Load(address string, user string, port int, dir string) {
	s.address = address
	s.user = user
	s.port = port
	s.dir = dir
}

func (s *Server) Dir() string {
	return fmt.Sprintf("/data/sites/%s", s.dir)
}

func (s *Server) Prefix() string {
	return fmt.Sprintf("[%s@%s]", s.user, s.address)
}

func (s *Server) Connect(privateKey string) error {
	if s.connOpened {
		log.Fatal("Warning: Client already connected")
		return nil
	}
	addr := fmt.Sprintf("%s:%d", s.address, s.port)
	replacePath, err := filepath.Abs(strings.Replace(privateKey, "~", os.Getenv("HOME"), 1))
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(replacePath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	s.conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	s.connOpened = true

	return nil
}

func (s *Server) Run(cmd string) error {
	sess, err := s.conn.NewSession()
	defer sess.Close()
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err = sess.RequestPty("vt220", 80, 40, modes); err != nil {
		return err
	}

	if s.stdin, err = sess.StdinPipe(); err != nil {
		return err
	}

	if s.stdout, err = sess.StdoutPipe(); err != nil {
		return err
	}

	if s.stderr, err = sess.StderrPipe(); err != nil {
		return err
	}

	if err = sess.Run(cmd); err != nil {
		return err
	}

	return nil
}

func (s *Server) CombinedOutput(cmd string) (out []byte, err error) {
	sess, err := s.conn.NewSession()
	defer sess.Close()
	if err != nil {
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err = sess.RequestPty("vt220", 80, 40, modes); err != nil {
		return
	}

	if s.stderr, err = sess.StderrPipe(); err != nil {
		return
	}

	out, err = sess.CombinedOutput(cmd)
	return
}

func (s *Server) Stdin() io.WriteCloser {
	return s.stdin
}

func (s *Server) Stdout() io.Reader {
	return s.stdout
}

func (s *Server) StdErr() io.Reader {
	return s.stderr
}

func (s *Server) Close() error {
	if !s.connOpened {
		log.Fatal("Warning: Trying to close the already closed connection")
		return nil
	}
	s.connOpened = false
	err := s.conn.Close()

	return err
}
