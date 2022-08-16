package user

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Task struct {
	Id   int
	Text string
	Time string
	Done bool
}

func NewStorage() []*Task {
	return []*Task{}
}

func NewConnectDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 password=1 dbname=API sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
