package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsantosalmeida/browser-chat/api/midleware"
	"github.com/vsantosalmeida/browser-chat/api/rest/handler"
	"github.com/vsantosalmeida/browser-chat/api/websocket"
	"github.com/vsantosalmeida/browser-chat/config"
	"github.com/vsantosalmeida/browser-chat/infrastructure/broker"
	"github.com/vsantosalmeida/browser-chat/infrastructure/repository"
	"github.com/vsantosalmeida/browser-chat/usecase/room"
	"github.com/vsantosalmeida/browser-chat/usecase/user"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func main() {
	config.InitLogging()

	db := config.InitDB()

	ch := config.InitRabbitMQ()

	// Setup User context
	userRepo := repository.NewUserMySQL(db)
	userSvc := user.NewService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Setup Room context
	roomRepo := repository.NewRoomMySQL(db)
	roomSvc := room.NewService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomSvc)

	// Setup WebSocket context
	rabbitMQ := broker.NewRabbitMQ(
		config.GetStingEnvVarOrPanic(config.ChatbotCommandOutputQueue), // read queue
		config.GetStingEnvVarOrPanic(config.ChatbotCommandInputQueue),  // write queue
		ch,
	)
	ctx, cancel := context.WithCancel(context.Background())
	wsServer := websocket.NewServer(roomSvc, rabbitMQ)

	go wsServer.Start(ctx)

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

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("unexpected server error: %v", err)
		}
	}()

	log.Info("server started")

	/// gracefully shutdown the server
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Info("server stopped")
}
