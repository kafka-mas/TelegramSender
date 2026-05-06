package main

import (
	"fmt"
	"log"

	"github.com/kafka-mas/TelegramSender/config"
	"github.com/kafka-mas/TelegramSender/database"
)

type Message interface {
	SendMessage(message *[]byte) error
}

type User interface {
	Create() error
	Delete() error
	Edit() error
	Select() error
}

type User_s struct {
	DBid   int
	TGid   int64
	TGname string
}

type DBase interface {
	Create() error
	Delete() error
	AddRecord(user_id int64, f_name string) (DBid int, err error)
	GetRecord(DBid int) (*database.DBUser, error)
	GetAllRecords() error
	RemoveRecord(DBid int) error
}

type Config interface {
	Read(filepath string) error
}

func main() {
	f := flags{}
	err := f.Parse()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	c := config.YamlConf{}
	err = c.Read("config.yaml")
	if err != nil {
		log.Fatalln("Error", err)
	}

	fmt.Println(c.Token)
	fmt.Println(c.Socks)
	fmt.Println()

	db := database.DBaseSQLite{Name: "./db.sqlite"}
	err = db.Create()
	if err != nil {
		log.Fatalln("Error", err)
	}

	id, err := db.AddRecord(123456789, "Alex")
	if err != nil {
		log.Println(err)
	}

	dUser, err := db.GetRecord(id)
	if err != nil {
		log.Fatalln(err)
	}
	myUser := User_s{
		DBid:   dUser.ID,
		TGid:   dUser.UserID,
		TGname: dUser.Fname,
	}

	fmt.Println("db ID:", myUser.DBid)
	fmt.Println("tg ID:", myUser.TGid)
	fmt.Println("name:", myUser.TGname)

	if err := db.Delete(); err != nil {
		log.Fatalln("Error", err)
	}
}
