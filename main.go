package main

import (
	"fmt"
	"log"

	"github.com/kafka-mas/TelegramSender/config"
	"github.com/kafka-mas/TelegramSender/database"
	"github.com/kafka-mas/TelegramSender/flagreader"
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

type Config interface {
	Read(filepath string) error
}

func main() {
	f := flagreader.Flags{}
	err := f.Parse()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	fmt.Println(f.IsAddUser)
	fmt.Println(f.Help)

	c := config.YamlConf{}
	err = c.Read("config.yaml")
	if err != nil {
		log.Fatalln("Error", err)
	}

	fmt.Println(c.Token)
	fmt.Println(c.Socks)
	fmt.Println()


	db := database.NewSQLite("./db.sqlite")
	err = db.Create()
	if err != nil {
		log.Fatalln("Error", err)
	}

	id, err := db.AddRecord(123456789, "Alex")
	if err != nil {
		log.Println(err)
	}

	dUser, err := db.GetRecord(id)
	if err != nil || dUser == nil {
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

	for range 3 {
		id, err = db.AddRecord(123456789, "Alex")
		if err != nil || id == -1 {
			log.Println(err)
		}
	}

	users, err := db.GetAllRecords()
	if err != nil {
		log.Println("Error get all records:", err)
	}
	fmt.Println()
	for _, u := range *users {
		fmt.Println(u.ID, u.UserID, u.Fname, u.ISdefault)
	}

	err = db.RemoveRecord(1)
	if err != nil {
		log.Println("Error remove record:", err)
	}

	err = db.SetDefault(2)
	if err != nil{
		log.Println(err)
	}

	users, err = db.GetAllRecords()
	if err != nil {
		log.Println("Error get all records:", err)
	}
	fmt.Println()
	for _, u := range *users {
		fmt.Println(u.ID, u.UserID, u.Fname, u.ISdefault)
	}

	if err := db.Delete(); err != nil {
		log.Fatalln("Error", err)
	}
}
