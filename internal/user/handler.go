package user

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	logger *logrus.Entry
	db     *sql.DB
}

func NewHandler(logger *logrus.Entry, postgres *sql.DB) *handler {
	return &handler{
		logger: logger,
		db:     postgres,
	}
}

func (h *handler) RegisterRouter(mux *mux.Router) {
	h.logger.Info("Регистрация обработчиков")
	mux.HandleFunc("/", h.IndexHandle).Methods("GET")
	//router.ServeFiles("/static/*filepath", http.Dir("static"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", h.IndexHandle)
	mux.HandleFunc("/delete/{id}", h.DeleteTask)
	mux.HandleFunc("/addTask", h.AddTask)
	mux.HandleFunc("/done/{id}", h.Done)
	mux.HandleFunc("/register", h.RegisterUser)
	mux.HandleFunc("/login", h.Login)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Домашняя страница")

	row, err := h.db.Query("select * from test order by id") // Соединение с БД
	if err != nil {
		panic(err)
	}

	var Tasks []Task // Новое хранилище для данных из БД
	for row.Next() {
		Tasks = append(Tasks, Task{})
		err := row.Scan(&Tasks[len(Tasks)-1].Id, &Tasks[len(Tasks)-1].Text, &Tasks[len(Tasks)-1].Time, &Tasks[len(Tasks)-1].Done)
		t := strings.NewReplacer("T", " ", "Z", "", "-", ".") // формат даты
		Tasks[len(Tasks)-1].Time = t.Replace(Tasks[len(Tasks)-1].Time)
		if err != nil {
			log.Fatal(err)
		}
	}
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(w, Tasks); err != nil {
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
		http.Redirect(w, r, "/", http.StatusNotAcceptable)
	}

	h.logger.Tracef("Пользователь зарегистрирован  %s : %s", name, password)
	//
	//rowAfected, err := result.RowsAffected()
	//if err != nil {
	//	h.logger.Warning(err)
	//}
	//
	//if rowAfected < 1 {
	//	h.logger.Warning(rowAfected)
	//}

	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	result := true

	row := h.db.QueryRow("select exists(select name, password from users where name = $1,$2)", name, password)
	row.Scan(&result)

	if !result {
		w.WriteHeader(http.StatusSeeOther)
	}

	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}
