package websocket

import (
	"net/http"
	"sync"

	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/pkg/auth"
	"github.com/vsantosalmeida/browser-chat/usecase/room"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrInvalidEventAction = errors.New("invalid event action")
)

type ClientList map[*Client]bool

// Server
type Server struct {
	clients     ClientList
	join        chan *Client
	leave       chan *Client
	handlers    map[string]EventHandler
	rooms       []*entity.Room
	roomUseCase room.UseCase
	mu          sync.RWMutex
}

// NewServer
func NewServer(roomUseCase room.UseCase) *Server {
	s := &Server{
		clients:     make(ClientList),
		join:        make(chan *Client),
		leave:       make(chan *Client),
		handlers:    initEventHandlers(),
		roomUseCase: roomUseCase,
	}

	rooms, err := s.roomUseCase.ListRooms()
	if err != nil {
		log.Fatalf("failed to load rooms: %v", err)
	}

	s.rooms = rooms

	return s
}

func (s *Server) Start() {
	log.Info("server started")

	for {
		select {

		case client := <-s.join:
			s.joinClient(client)

		case client := <-s.leave:
			s.leaveClient(client)
		}
	}
}

// ServeWS
func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request) {
	userCtxValue := r.Context().Value(auth.UserContextKey)
	if userCtxValue == nil {
		log.Error("unauthorized connection")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := userCtxValue.(entity.AuthenticatedUser)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(conn, s, user.GetUsername(), user.GetId())

	go client.readMessages()
	go client.writeMessages()

	s.join <- client
}

func (s *Server) joinClient(client *Client) {
	s.clients[client] = true
	log.WithField("UserID", client.ID).Info("user connected")
}

func (s *Server) leaveClient(client *Client) {
	if _, ok := s.clients[client]; ok {
		client.conn.Close()
		delete(s.clients, client)
		log.WithField("UserID", client.ID).Info("user disconnected")
	}
}

// routeEvent
func (s *Server) routeEvent(event Event, c *Client) error {
	if handler, ok := s.handlers[event.Action]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		log.WithField("Action", event.Action).Error("invalid event action")
		return ErrInvalidEventAction
	}
}

func (s *Server) isValidRoom(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if found := findRoom(s.rooms, id); found {
		return true
	} else {
		// if room isn't found server will try to retrieve from DB as a last chance
		rooms, err := s.roomUseCase.ListRooms()
		if err != nil {
			return false
		}

		s.rooms = rooms

		return findRoom(s.rooms, id)
	}
}

func findRoom(rooms []*entity.Room, id int) bool {
	for _, r := range rooms {
		if r.ID == id {
			return true
		}
	}

	log.WithField("RoomID", id).Warn("room not found")
	return false
}

func initEventHandlers() map[string]EventHandler {
	handlers := map[string]EventHandler{
		SendMessageAction: SendMessageHandler,
		JoinRoomAction:    ChatRoomHandler,
	}

	return handlers
}
