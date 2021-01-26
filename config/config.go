package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Server Host `yaml:server`
}

type Host struct {
	Address string `yaml:"address"`
	User    string `yaml:"user"`
	Port    int    `yaml:"port"`
	Dir     string `yaml:"dir"`
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
