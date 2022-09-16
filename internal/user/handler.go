package user

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
var tmpl = template.Must(template.ParseFiles("index.html", "login.html", "register.html"))

type handler struct {
	logger *logrus.Entry
	db     *sql.DB
	user   *User
}

func NewHandler(logger *logrus.Entry, postgres *sql.DB, user *User) *handler {
	return &handler{
		logger: logger,
		db:     postgres,
		user:   user,
	}
}

func (h *handler) RegisterRouter(mux *mux.Router) {
	h.logger.Info("Регистрация обработчиков")

	mux.HandleFunc("/", h.IndexHandle)
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/delete/{id}", CheckCookie(h.DeleteTask))
	mux.HandleFunc("/addTask/", CheckCookie(h.AddTask))
	mux.HandleFunc("/done/{id}", CheckCookie(h.Done))
	mux.HandleFunc("/register", h.RegisterUser)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/logout", h.Logout)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request) {
	if h.user.Entry {
		query := fmt.Sprintf("select * from %s order by id", h.user.Name)

		row, err := h.db.Query(query)
		if err != nil {
			h.logger.Error(err)
		}

		result, _ := h.db.Exec(query)
		numOfColumns, _ := result.RowsAffected()

		h.user.Tasks = make([]Task, numOfColumns)

		for i := 0; row.Next(); i++ {
			err := row.Scan(&h.user.Tasks[i].Id, &h.user.Tasks[i].Text, &h.user.Tasks[i].Time, &h.user.Tasks[i].Done)
			t := strings.NewReplacer("T", " ", "Z", "", "-", ".") // формат даты
			h.user.Tasks[i].Time = t.Replace(h.user.Tasks[i].Time)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err := tmpl.Execute(w, h.user); err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		//h.user.SaveJSON()

	} else {
		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	table := fmt.Sprintf("delete from %s where id = %s", h.user.Name, vars["id"])
	_, err := h.db.Exec(table)
	if err != nil {
		h.logger.Error(err)
	}

	h.logger.Warnf("Удаление записи id : %s", vars["id"])
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	task := r.FormValue("text")

	table := fmt.Sprintf("insert into %s(task) values ('%s')", h.user.Name, task)

	h.logger.Warning(table)

	_, err := h.db.Exec(table)
	if err != nil {
		h.logger.Error(err)
	}
	h.logger.Infof("Добавление записи в БД : '%s'", task)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return
}

func (h *handler) Done(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	table := fmt.Sprintf("update %s set done = not done where id = %s", h.user.Name, vars["id"])
	_, err := h.db.Exec(table)
	if err != nil {
		h.logger.Error(err)
	}
	h.logger.Infof("Измениние состояния id = : %v", vars["id"])
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	name := r.FormValue("name")
	password := r.FormValue("password")

	err := tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		h.logger.Error(err)
	}

	h.user.Name = name
	h.user.Password = password
	//create table if not exist
	table := fmt.Sprintf("create table %s (id serial,task text,time timestamp default now(), done bool default true)", name)
	_, err = h.db.Exec(table)

	h.logger.Errorf("Create table : %s", name)

	_, err = h.db.Exec("insert into users(name, password) values ($1,$2)", name, password)
	if err != nil {
		h.logger.Warning("Ошибка регистрации\n", err)
		http.Redirect(w, r, "/", 303)
	}

	h.logger.Tracef("Пользователь зарегистрирован  %s : %s", name, password)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	name := r.FormValue("name")
	password := r.FormValue("password")

	row := h.db.QueryRow("select exists(select name, password from users where name =$1 and password=$2)", name, password)
	if err := row.Scan(&h.user.Entry); err != nil {
		h.logger.Warning(err)
	}

	if h.user.Entry == true {
		session, _ := store.Get(r, "cookie-name")
		session.Values["auth"] = true
		if err := session.Save(r, w); err != nil {
			h.logger.Error(err)
		}
		h.user.Name = name
		h.user.Password = password

		//file, err := os.OpenFile(h.user.DBName, os.O_RDWR|os.O_APPEND, 0644)
		//if err != nil {
		//	h.logger.Error(err)
		//}
		//encoder := json.NewDecoder(file)
		//err = encoder.Decode(&h.user.Tasks)
		//if err != nil {
		//	h.logger.Error(err, h.user)
		//}
		//
		//defer file.Close()
		http.Redirect(w, r, "/", 303)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		h.logger.Error(err)
	}
	session.Values["auth"] = false
	if err := session.Save(r, w); err != nil {
		h.logger.Error(err)
	}
	h.user.Entry = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
