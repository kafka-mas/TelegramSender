package yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlConf struct {
	Token string `yaml:"token"`
	Socks struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
	} `yaml:"socks"`
}

type ConfPath struct{ Filename string }

func (c ConfPath) Read() (*YamlConf, error) {
	data, err := os.ReadFile(c.Filename)
	if err != nil {
		err = fmt.Errorf("Error open file: %v", err)
		return nil, err
	}
	conf := &YamlConf{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		err = fmt.Errorf("Error parse YAML: %v", err)
		return nil, err
	}

	return conf, nil
}
