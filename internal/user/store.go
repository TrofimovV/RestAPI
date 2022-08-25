package user

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type User struct {
	Tasks    []Task
	Name     string
	Password string
	Entry    bool
}

type Task struct {
	Id   int
	Text string
	Time string
	Done bool
}

func NewConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "password=1 dbname=API")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewUser() *User {
	return &User{}
}
