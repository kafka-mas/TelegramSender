package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kafka-mas/TelegramSender/config"
	_ "github.com/kafka-mas/TelegramSender/database"
	"github.com/kafka-mas/TelegramSender/flagreader"
	"github.com/kafka-mas/TelegramSender/telegram"
)

type User interface {
	Create() error
	Delete() error
	Edit() error
	Select() error
}

type User_s struct {
	DBid      int
	TGid      int64
	TGname    string
	IsDefault bool
}

/*
			a 	d 	delete 	list 	s 	h 	id 	username	version
	a		1	1	0		0		0	0	0	0			0
	d			1	1		0		1	0	0	0			0
	delete			1		0		0	0	1	1			0
	list					1		0	0	0	0			0
	s								1	0	1	1			0
	h									1	0	0			0
	id										1	0			0
	username									1			0
	version													1
*/

func main() {
	f := flagreader.Flags{}
	err := f.Parse()
	if err != nil {
		log.Fatalln("error:", err)
	}

	confPath, err := confSelect()
	if err != nil {
		log.Fatalln("error:", err)
	}
	c := config.Yaml(confPath)
	conf, err := c.Read()
	if err != nil {
		log.Fatalln("error", err)
	}

	token := conf.Token
	socks := struct {
		addr string
		port int
		user string
		pass string
	}{
		addr: conf.Socks.Address,
		port: conf.Socks.Port,
		user: conf.Socks.User,
		pass: conf.Socks.Pass,
	}

	bot := telegram.NewBot(token, telegram.WithSocksProxy(socks.addr, socks.port))

	// err = bot.SendFile(f.UserId, "./telegram/telegram.go")
	// if err != nil {
	// 	log.Println(err)
	// }
	_, err = bot.ReceiveFile(15)
	if err != nil {
		log.Fatalln("Error:", err)
	}



}

func confSelect() (string, error) {
	home, _ := os.UserHomeDir()
	paths := []string{
		"/etc/telegram_sender/config.yaml",
		filepath.Join(home, "/.config/telegram_sender/config.yaml"),
		"./config.yaml",
	}
	fmt.Println(paths)
	var confPath string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				return "", fmt.Errorf("error open config file: %w", err)
			}
		}
		if !info.IsDir() {
			confPath = path
		}
	}
	if confPath == "" {
		return "", fmt.Errorf("error, config file not found")
	}

	return confPath, nil
}
