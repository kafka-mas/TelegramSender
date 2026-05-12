package config

import "github.com/kafka-mas/TelegramSender/config/yaml"

type Config interface {
	Read() (*yaml.YamlConf, error)
}

func Yaml(filename string) Config { return yaml.ConfPath{Filename: filename} }
