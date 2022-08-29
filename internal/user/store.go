package user

import (
	"RestAPI/pkg/logging"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"os"
)

type User struct {
	Tasks    []Task `json:"tasks,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"-"`
	Entry    bool   `json:"entry"`
}

type Task struct {
	Id   int    `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	Time string `json:"time,omitempty"`
	Done bool   `json:"done,omitempty"`
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

func (u *User) SaveJSON() {
	logger := logging.GetLogger()

	marshal, err := json.Marshal(&u)
	if err != nil {
		logger.Error(err)
	}

	file, err := os.Create(u.Name)
	if err != nil {
		logger.Error(err)
	}

	defer file.Close()

	_, err = file.Write(marshal)
	if err != nil {
		logger.Error(err)
	}
}
