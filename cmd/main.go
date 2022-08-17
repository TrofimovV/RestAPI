package main

import (
	"RestAPI/internal/user"
	"RestAPI/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
)

func main() {

	log := logging.GetLogger()

	router := httprouter.New()

	handler := user.NewHandler(log)

	handler.RegisterRouter(router)

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

	logging.GetLogger().Debugf("Сервер слушает порт : 8080")

	server.Serve(listener)

}
