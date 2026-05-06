package main

import (
	"fmt"
	"flag"
)

type flags struct {
	isAddUser    bool
	isDefault    bool
	isDeleteUser bool
	isListUser   bool
	isString     string
	help         bool
	userId       int64
	userName     string
	version      bool
}

func (f *flags) Parse() error {
	flag.BoolVar(&f.isAddUser, "a", f.isAddUser, "Add new user to database")
	flag.BoolVar(&f.isDefault, "d", f.isDefault, "Send message to default user")
	flag.BoolVar(&f.isDeleteUser, "delete", f.isDeleteUser, "Delete user from database")
	flag.BoolVar(&f.isListUser, "list", f.isListUser, "List saved users")
	flag.StringVar(&f.isString, "s", f.isString, "Send string")
	flag.BoolVar(&f.help, "h", f.help, "Show help")
	flag.Int64Var(&f.userId, "id", f.userId, "Send specified user by Telegram ID")
	flag.StringVar(&f.userName, "username", f.userName, "Send to specified username")
	flag.BoolVar(&f.version, "v", f.version, "Show version")
	flag.Parse()

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

	if f.isAddUser == true {
		if f.isDeleteUser == true ||
			f.isListUser == true ||
			f.isString != "" ||
			f.help == true ||
			f.userId != 0 ||
			f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.isDefault == true {
		if f.isListUser == true ||
			f.isString != "" ||
			f.help == true ||
			f.userId != 0 ||
			f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.isDeleteUser == true {
		if f.isListUser == true ||
			f.isString != "" ||
			f.help == true ||
			f.userId != 0 ||
			f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.isListUser == true {
		if f.isString != "" ||
			f.help == true ||
			f.userId != 0 ||
			f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.isString != "" {
		if f.help == true ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.help == true {
		if f.userId != 0 ||
			f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.userId != 0 {
		if f.userName != "" ||
			f.version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.userName != "" && f.version == true {
		return fmt.Errorf("incompatible parameters")
	}

	return nil
}
