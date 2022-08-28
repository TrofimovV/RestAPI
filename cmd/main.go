package main

import (
	"RestAPI/internal/user"
	"RestAPI/pkg/logging"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Инициализация логгера")

	router := mux.NewRouter()

	db, err := user.NewConnectDB()
	if err != nil {
		logger.Fatal()
	}

	u := user.NewUser()

	handler := user.NewHandler(logger, db, u)

	handler.RegisterRouter(router)

	//trying to get json

	user.DecodeJSON(u)

	if err != start(router) {
		logger.Fatal(err)
	}
}

func start(router *mux.Router) error {
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
