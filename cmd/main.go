package main

import (
	"RestAPI/internal/user"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()

	handler := user.NewHandler()

	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	http.ListenAndServe(":8080", router)
	//listener, err := net.Listen("tcp", ":1234")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//server := &http.Server{
	//	Handler:      router,
	//	ReadTimeout:  10,
	//	WriteTimeout: 10,
	//}
	//log.Fatal(server.Serve(listener))
}
