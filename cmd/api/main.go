package main

import (
	"net/http"
	"time"

	"github.com/vsantosalmeida/browser-chat/api/midleware"
	"github.com/vsantosalmeida/browser-chat/api/rest/handler"
	"github.com/vsantosalmeida/browser-chat/api/websocket"
	"github.com/vsantosalmeida/browser-chat/config"
	"github.com/vsantosalmeida/browser-chat/infrastructure/repository"
	"github.com/vsantosalmeida/browser-chat/usecase/room"
	"github.com/vsantosalmeida/browser-chat/usecase/user"

	"github.com/apex/log"
	"github.com/gorilla/mux"
)

func main() {
	config.InitLogging()

	db := config.InitDB()

	// Setup User context
	userRepo := repository.NewUserMySQL(db)
	userSvc := user.NewService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Setup Room context
	roomRepo := repository.NewRoomMySQL(db)
	roomSvc := room.NewService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomSvc)

	// Setup WebSocket context
	wsServer := websocket.NewServer(roomSvc)
	go wsServer.Start()

	// Setup HTTP handlers
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.HandleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", userHandler.HandleListUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/login", userHandler.HandleLogin).Methods(http.MethodPost)

	r.HandleFunc("/rooms", roomHandler.HandleCreateRoom).Methods(http.MethodPost)
	r.HandleFunc("/rooms/{id}/messages", roomHandler.HandleListMessages).Methods(http.MethodGet)
	r.HandleFunc("/rooms", roomHandler.HandleListRooms).Methods(http.MethodGet)

	r.HandleFunc("/ws", midleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		wsServer.ServeWS(w, r)
	}))

	r.Use(midleware.Cors)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("closing server")
	}
}
