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
	tmpl, err := template.ParseFiles("index.html", "login.html", "register.html")
	if err != nil {
		panic(err)
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", h.user); err != nil {
		panic(err)
	}
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	h.logger.Warnf("Удаление записи id : %s", vars["id"])
	_, err := h.db.Exec("delete from test where id = $1", vars["id"])
	if err != nil {
		h.logger.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("text")

	_, err := h.db.Exec("insert into test(task) values ($1)", task)
	if err != nil {
		h.logger.Fatal(err)
	}
	h.logger.Infof("Добавление записи в БД : '%s'", task)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) Done(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.logger.Infof("Измениние состояния id = : %v", vars["id"])
	_, err := h.db.Exec("update test set done = not done where id = $1", vars["id"])
	if err != nil {
		h.logger.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	/*result,*/
	_, err := h.db.Exec("insert into users(name, password) values ($1,$2)", name, password)
	if err != nil {
		h.logger.Warning("Ошибка регистрации\n", err)
		http.Redirect(w, r, "/login", http.StatusNotAcceptable)
	}

	h.logger.Tracef("Пользователь зарегистрирован  %s : %s", name, password)

	defer http.Redirect(w, r, "/", http.StatusOK)
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	row := h.db.QueryRow("select exists(select name, password from users where name =$1 and password=$2)", name, password)
	if err := row.Scan(&h.user.Entry); err != nil {
		h.logger.Warning(err)
	}

	if h.user.Entry == true {
		session, _ := store.Get(r, "cookie-name")
		session.Values["auth"] = true
		session.Save(r, w)
		h.user.Name = name
		h.user.Password = password
		http.Redirect(w, r, "/", 200)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["auth"] = false
	session.Save(r, w)
	h.user.Entry = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
