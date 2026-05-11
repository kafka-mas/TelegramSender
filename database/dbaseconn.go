package database

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBUser struct {
	ID        int
	UserID    int64
	Fname     string
	ISdefault bool
}

type DB interface {
	Create() error
	Delete() error
	AddRecord(user_id int64, f_name string) (DBid int, err error)
	GetRecord(DBid int) (*DBUser, error)
	GetAllRecords() (*[]DBUser, error)
	RemoveRecord(DBid int) error
	SetDefault(DBid int) error
	RemoveDefault() error
}

func NewSQLite(name string) DB { return sqliteDB{name: name} }

type sqliteDB struct{ name string }

func (d sqliteDB) Create() error {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	createTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER UNIQUE,
            f_name TEXT,
            is_default INTEGER NOT NULL DEFAULT 0
        );`
	if _, err = db.Exec(createTable); err != nil {
		return fmt.Errorf("error create table: %v", err)
	}

	// Создание частичного уникального индекса
	createIndex := `
        CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_default 
        ON users(is_default) WHERE is_default = 1;`
	if _, err = db.Exec(createIndex); err != nil {
		return fmt.Errorf("error create index: %v", err)
	}
	log.Println("The database was created successfully (or was created earlier)")

	return nil
}

func (d sqliteDB) Delete() error {
	err := os.Remove(d.name)
	if err != nil {
		return fmt.Errorf("error delete database %v: %v", d.name, err)
	}
	log.Println("Database succesfully deleted.")
	return nil
}

func (d sqliteDB) AddRecord(user_id int64, f_name string) (DBid int, err error) {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return -1, fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	// char *sql = "INSERT INTO users (user_id, f_name, is_default) VALUES (?, ?, ?);";
	result, err := db.Exec("insert into users (user_id, f_name, is_default) values ($1, $2, 0)", user_id, f_name)
	if err != nil {
		return -1, fmt.Errorf("error add record %v in database %v: %v", f_name, d.name, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get last insert id: %v", err)
	}

	if lastID > math.MaxInt {
		return -1, fmt.Errorf("last insert id %d exceeds int max value (%d)", lastID, math.MaxInt)
	}

	DBid = int(lastID)
	return
}

func (d sqliteDB) GetRecord(DBid int) (*DBUser, error) {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return nil, fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	row := db.QueryRow("select id, user_id, f_name from users where id = $1", DBid)
	u := DBUser{}
	err = row.Scan(&u.ID, &u.UserID, &u.Fname)
	if err != nil {
		return nil, fmt.Errorf("error get value by id %d: %v", DBid, err)
	}

	return &u, nil
}

func (d sqliteDB) GetAllRecords() (*[]DBUser, error) {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return nil, fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []DBUser{}
	for rows.Next() {
		u := DBUser{}
		rows.Scan(&u.ID, &u.UserID, &u.Fname, &u.ISdefault)
		users = append(users, u)
	}

	return &users, nil
}

func (d sqliteDB) RemoveRecord(DBid int) error {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("delete from users where id = $1", DBid)
	if err != nil {
		return fmt.Errorf("error delete record: %v", err)
	}
	fmt.Printf("Record with ID %d succesfully removed.\n", DBid)
	return nil
}

// SetDefault(DBid int) error

func (d sqliteDB) SetDefault(DBid int) error {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("update users set is_default = 0")
	if err != nil {
		return err
	}
	_, err = db.Exec("update users set is_default = 1 where id = $1", DBid)
	if err != nil {
		return err
	}

	return nil
}

func (d sqliteDB) RemoveDefault() error {
	db, err := sql.Open("sqlite3", d.name)
	if err != nil {
		return fmt.Errorf("error open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("update users set is_default = 0")
	if err != nil {
		return err
	}

	return nil
}
