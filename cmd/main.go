package main

import (
	"RestAPI/internal/user"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
)

func main() {
	router := httprouter.New()

	handler := user.NewHandler()

	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler: router,
	}
	log.Fatal(server.Serve(listener))
}
