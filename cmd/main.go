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
	logger := logging.GetLogger()

	logger.Info("Инициализация логгера")

	router := httprouter.New()

	db, err := user.NewConnectDB()
	if err != nil {
		logger.Fatal()
	}

	handler := user.NewHandler(logger, db)

	handler.RegisterRouter(router)

	//user.RegisterUser("1", "1")

	if err != start(router) {
		logger.Fatal(err)
	}
}

func start(router *httprouter.Router) error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler: router,
	}

	logging.GetLogger().Debug("Сервер слушает порт : 8080")

	return server.Serve(listener)
}
