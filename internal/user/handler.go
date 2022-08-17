package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	log *logrus.Entry
}

func NewHandler(logger *logrus.Entry) *handler {
	return &handler{
		log: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	h.log.Trace("11111111")
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/", h.IndexHandle)
	router.GET("/delete/:uuid", h.DeleteTask)
	router.GET("/addTask/", h.AddTask)
	router.GET("/done/:uuid", h.Done)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.log.Info("fddjsfdj")
	row, err := NewConnectDB().Query("select * from test order by id") // Соединение с БД
	if err != nil {
		panic(err)
	}
	var Tasks []Task // Новое хранилище для бд
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
	_, err := NewConnectDB().Exec(fmt.Sprintf("delete from test where id = %v", p))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Удаление поля")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	task := r.FormValue("text")
	_, err := NewConnectDB().Exec("insert into test(task) values ($1)", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Добавление поля ")
	defer http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *handler) Done(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p := params.ByName("uuid")
	_, err := NewConnectDB().Exec("update test set done = not done where id = $1", p)
	if err != nil {
		log.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}
