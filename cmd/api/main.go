package main

import (
	"github.com/vsantosalmeida/browser-chat/api/midleware"
	"log"
	"net/http"
	"time"

	"github.com/vsantosalmeida/browser-chat/api/rest/handler"
	"github.com/vsantosalmeida/browser-chat/config"
	"github.com/vsantosalmeida/browser-chat/infrastructure/repository"
	"github.com/vsantosalmeida/browser-chat/usecase/user"

	"github.com/gorilla/mux"
)

func main() {
	db := config.InitDB()

	// Setup User context
	userRepo := repository.NewUserMySQL(db)
	userSvc := user.NewService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Setup HTTP handlers
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.HandleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", userHandler.HandleListUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/login", userHandler.HandleLogin).Methods(http.MethodPost)
	r.Use(midleware.Cors)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
