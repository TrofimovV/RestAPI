package user

import (
	"RestAPI/configs"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type User struct {
	Tasks    []Task `json:"-"`
	Name     string `json:"-"`
	Password string `json:"-"`
	Entry    bool   `json:"-"`
}

type Task struct {
	Id   int    `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	Time string `json:"time,omitempty"`
	Done bool   `json:"done,omitempty"`
}

func NewConnectDB(logger *logrus.Entry, cfg *configs.ConfigDatabase) (*sql.DB, error) {
	//dataConfig := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
	dataConfig := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	//logging with "pkg/logging"
	logger.Debugf("POSTGRES_DB = %s POSTGRES_USER = %s", cfg.DBName, cfg.Username)

	db, err := sql.Open("postgres", dataConfig)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error : %v", err)
	}
	return db, nil
}

func NewUser() *User {
	return &User{}

}

//func (u *User) SaveJSON() {
//	logger := logging.GetLogger()
//
//	marshal, err := json.MarshalIndent(&u.Tasks, "", "")
//	if err != nil {
//		logger.Error(err)
//	}
//
//	file, err := os.Create(u.DBName)
//	if err != nil {
//		logger.Errorf("имя пользователя: %s : %v", u.DBName, err)
//	}
//
//	defer file.Close()
//
//	_, err = file.Write(marshal)
//	if err != nil {
//		logger.Error(err)
//	}
//	logger.Info("Save JSON ")
//}
