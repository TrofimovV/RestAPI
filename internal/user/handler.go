package user

import (
	"database/sql"
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

func (h *handler) IndexHandle(w http.ResponseWriter, _ *http.Request) {
	row, err := h.db.Query("select * from test order by id")

	if err != nil {
		panic(err)
	}

	result, _ := h.db.Exec("select * from test order by id")
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

	h.user.SaveJSON()

}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	h.logger.Warnf("Удаление записи id : %s", vars["id"])
	_, err := h.db.Exec("delete from test where id = $1", vars["id"])
	if err != nil {
		h.logger.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("text")

	_, err := h.db.Exec("insert into test(task) values ($1)", task)
	if err != nil {
		h.logger.Fatal(err)
	}
	h.logger.Infof("Добавление записи в БД : '%s'", task)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *handler) Done(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.logger.Infof("Измениние состояния id = : %v", vars["id"])
	_, err := h.db.Exec("update test set done = not done where id = $1", vars["id"])
	if err != nil {
		h.logger.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	err := tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		h.logger.Error(err)
	}

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

		//file, err := os.OpenFile(h.user.Name, os.O_RDWR|os.O_APPEND, 0644)
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
