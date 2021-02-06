package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Server     Host `yaml:server`
	Repository Repo `yaml:repository`
}

type Host struct {
	Address string `yaml:"address"`
	User    string `yaml:"user"`
	Port    int    `yaml:"port"`
	Dir     string `yaml:"dir"`
}

type Repo struct {
	Url    string `yaml:"url"`
	Branch string `yaml:"branch"`
	Tag    string `yaml:"tag"`
}

func (c *Configuration) ReadFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	errY := yaml.Unmarshal(file, &c)
	if errY != nil {
		return errY
	}

	return nil
}
