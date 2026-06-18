package flagreader

import (
	"fmt"
	"flag"
)

type Flags struct {
	IsAddUser    bool
	IsDefault    bool
	IsDeleteUser bool
	IsListUser   bool
	IsString     string
	Help         bool
	UserId       int64
	UserName     string
	Version      bool
}

func (f *Flags) Parse() error {
	flag.BoolVar(&f.IsAddUser, "a", f.IsAddUser, "Add new user to database")
	flag.BoolVar(&f.IsDefault, "d", f.IsDefault, "Send message to default user")
	flag.BoolVar(&f.IsDeleteUser, "delete", f.IsDeleteUser, "Delete user from database")
	flag.BoolVar(&f.IsListUser, "list", f.IsListUser, "List saved users")
	flag.StringVar(&f.IsString, "s", f.IsString, "Send string")
	flag.BoolVar(&f.Help, "h", f.Help, "Show help")
	flag.Int64Var(&f.UserId, "id", f.UserId, "Send specified user by Telegram ID")
	flag.StringVar(&f.UserName, "username", f.UserName, "Send to specified username")
	flag.BoolVar(&f.Version, "v", f.Version, "Show version")
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

	if f.IsAddUser == true {
		if f.IsDeleteUser == true ||
			f.IsListUser == true ||
			f.IsString != "" ||
			f.Help == true ||
			f.UserId != 0 ||
			f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.IsDefault == true {
		if f.IsListUser == true ||
			f.IsString != "" ||
			f.Help == true ||
			f.UserId != 0 ||
			f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.IsDeleteUser == true {
		if f.IsListUser == true ||
			f.IsString != "" ||
			f.Help == true ||
			f.UserId != 0 ||
			f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.IsListUser == true {
		if f.IsString != "" ||
			f.Help == true ||
			f.UserId != 0 ||
			f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.IsString != "" {
		if f.Help == true ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.Help == true {
		if f.UserId != 0 ||
			f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.UserId != 0 {
		if f.UserName != "" ||
			f.Version == true {
			return fmt.Errorf("incompatible parameters")
		}
	}

	if f.UserName != "" && f.Version == true {
		return fmt.Errorf("incompatible parameters")
	}

	return nil
}
