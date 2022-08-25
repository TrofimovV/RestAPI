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

var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))
var result bool

type handler struct {
	logger *logrus.Entry
	db     *sql.DB
	user   *User
}

func NewHandler(logger *logrus.Entry, postgres *sql.DB, u *User) *handler {
	return &handler{
		logger: logger,
		db:     postgres,
		user:   u,
	}
}

func (h *handler) RegisterRouter(mux *mux.Router) {
	h.logger.Info("Регистрация обработчиков")

	mux.HandleFunc("/", h.IndexHandle)
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/delete/{id}", h.DeleteTask)
	mux.HandleFunc("/addTask/", h.AddTask)
	mux.HandleFunc("/done/{id}", h.Done)
	mux.HandleFunc("/register", h.RegisterUser)
	mux.HandleFunc("/login", h.Login)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request) {
	//session, _ := store.Get(r, "cookie-name")
	//
	//if auth, ok := session.Values["auth"].(bool); !ok || !auth {
	//	http.Error(w, "Forbidden", http.StatusForbidden)
	//	return
	//}
	h.logger.Info("Домашняя страница")

	row, err := h.db.Query("select * from test order by id") // Соединение с БД
	if err != nil {
		panic(err)
	}

	u := User{} // Новое хранилище для данных из БД

	for row.Next() {
		u.Tasks = append(u.Tasks, Task{})
		err := row.Scan(&u.Tasks[len(u.Tasks)-1].Id, &u.Tasks[len(u.Tasks)-1].Text, &u.Tasks[len(u.Tasks)-1].Time, &u.Tasks[len(u.Tasks)-1].Done)
		t := strings.NewReplacer("T", " ", "Z", "", "-", ".") // формат даты
		u.Tasks[len(u.Tasks)-1].Time = t.Replace(u.Tasks[len(u.Tasks)-1].Time)
		if err != nil {
			log.Fatal(err)
		}
	}
	tmpl, err := template.ParseFiles("index.html", "login.html", "register.html")
	if err != nil {
		panic(err)
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", u); err != nil {
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
	if err := row.Scan(&result); err != nil {
		h.logger.Warning(err)
	}

	if result {
		session, _ := store.Get(r, "cookie-name")
		session.Values["auth"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", 200)
	}

	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}
