package user

import (
	"RestAPI/pkg/logging"
	"database/sql"
	"encoding/json"
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

func DecodeJSON(u *User) {
	logger := logging.GetLogger()
	marshal, err := json.Marshal(u)
	if err != nil {
		return
	}
	logger.Warning(marshal)
}

//TODO GET JSOM FROM USER
func SaveTable() {

}
