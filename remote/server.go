package remote

import (
	"fmt"
	"golang.org/x/crypto/ssh"
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
}

func (s *Server) Load(address string, user string, port int, dir string) {
	s.address = address
	s.user = user
	s.port = port
	s.dir = dir
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

func (s *Server) Close() error {
	if !s.connOpened {
		log.Fatal("Warning: Trying to close the already closed connection")
		return nil
	}
	s.connOpened = false
	err := s.conn.Close()

	return err
}
