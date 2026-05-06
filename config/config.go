package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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

func (c *YamlConf) Read(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("Error open file: %v", err)
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		err = fmt.Errorf("Error parse YAML: %v", err)
		return err
	}

	return nil
}
