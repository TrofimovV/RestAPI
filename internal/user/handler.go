package user

import (
	"RestAPI/internal/handlers"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type handler struct {
}

func New() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/", h.IndexHandle)
	router.DELETE("/delete/", h.DeleteTask)
	router.POST("/addTask/", h.AddTask)
	router.PATCH("/done/", h.Done)
}

func (h *handler) IndexHandle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	row, err := user.NewConnectDB().Query("select * from test order by id") // Соединение с БД
	if err != nil {
		panic(err)
	}
	var Tasks []user.Task // Новое хранилище для бд
	for row.Next() {
		Tasks = append(Tasks, user.Task{})
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
	IdTask := r.URL.Path
	IdTask = strings.TrimLeft(IdTask, "/delete/")
	_, err := user.NewConnectDB().Exec(fmt.Sprintf("delete from test where id = %v", IdTask))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Удаление поля")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	task := r.FormValue("text")
	_, err := user.NewConnectDB().Exec("insert into test(task) values ($1)", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Добавление поля ")
	defer http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *handler) Done(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	done := r.URL.Path
	done = strings.TrimLeft(done, "/done/")
	_, err := user.NewConnectDB().Exec("update test set done = not done where id = $1", done)
	if err != nil {
		log.Fatal(err)
	}
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
}