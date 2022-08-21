package user

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
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

func (h *handler) RegisterRouter(router *httprouter.Router) {
	h.logger.Info("Регистрация обработчиков")

	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/", h.IndexHandle)
	router.GET("/delete/:uuid", h.DeleteTask)
	router.GET("/addTask/", h.AddTask)
	router.GET("/done/:uuid", h.Done)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p := params.ByName("uuid")
	h.logger.Warnf("Удаление записи id : %s", p)
	_, err := h.db.Exec("delete from test where id = $1", p)
	if err != nil {
		h.logger.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	task := r.FormValue("text")
	_, err := h.db.Exec("insert into test(task) values ($1)", task)
	if err != nil {
		h.logger.Fatal(err)
	}
	h.logger.Infof("Добавление записи в БД : '%s'", task)
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) Done(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p := params.ByName("uuid")
	h.logger.Infof("Измениние состояния id = : %v", p)
	_, err := h.db.Exec("update test set done = not done where id = $1", p)
	if err != nil {
		h.logger.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//name := r.FormValue("name")
}
