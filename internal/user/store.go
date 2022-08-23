package user

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Task struct {
	Id   int
	Text string
	Time string
	Done bool
}

//todo сделать приветсвенное окно
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
