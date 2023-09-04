package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

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

	// ErrInvalidEventAction invalid event action
	ErrInvalidEventAction = errors.New("invalid event action")
)

// ClientList holds the current connected Clients with the Server.
type ClientList map[*Client]bool

// Server handle the websocket connection between Clients and events.
type Server struct {
	clients     ClientList
	join        chan *Client
	leave       chan *Client
	handlers    map[string]EventHandler
	rooms       []*entity.Room
	roomUseCase room.UseCase
	broker      Broker
	mu          sync.RWMutex
}

// CommandOutput result of executed command from chatbot.
type CommandOutput struct {
	RoomID  int    `json:"roomID"`
	From    string `json:"from"`
	Message string `json:"message"`
}

// NewServer Server builder.
func NewServer(roomUseCase room.UseCase, broker Broker) *Server {
	s := &Server{
		clients:     make(ClientList),
		join:        make(chan *Client),
		leave:       make(chan *Client),
		handlers:    initEventHandlers(),
		roomUseCase: roomUseCase,
		broker:      broker,
	}

	rooms, err := s.roomUseCase.ListRooms()
	if err != nil {
		log.Fatalf("failed to load rooms: %v", err)
	}

	s.rooms = rooms

	return s
}

// Start loop to receive Client connections or disconnections.
func (s *Server) Start(ctx context.Context) {
	go s.listenChatbotMessages(ctx)

	for {
		if ctx.Err() == context.Canceled {
			log.Warn("context canceled")
			break
		}

		select {
		case client := <-s.join:
			s.joinClient(client)

		case client := <-s.leave:
			s.leaveClient(client)
		}
	}
}

// ServeWS handle the websocket connections with an authenticated Client and starts the go routines
// to listen for read and write events.
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

// joinClient adds a connected Client to the Server.
func (s *Server) joinClient(client *Client) {
	s.clients[client] = true
	log.WithField("UserID", client.ID).Info("user connected")
}

// leaveClient disconnects a Client from the Server.
func (s *Server) leaveClient(client *Client) {
	if _, ok := s.clients[client]; ok {
		client.conn.Close()
		delete(s.clients, client)
		log.WithField("UserID", client.ID).Info("user disconnected")
	}
}

// routeEvent find the EventHandler for the respective event and process it.
// it throws an error if the EventHandler is not found.
func (s *Server) routeEvent(event Event, c *Client) error {
	if handler, ok := s.handlers[event.Action]; ok {
		return handler(event, c)
	}

	log.WithField("Action", event.Action).Error("invalid event action")
	return ErrInvalidEventAction
}

// isValidRoom checks if the given room ID exists in the Server memory.
// will try to retrieve from DB if this chat room isn't in the memory.
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

// listenChatbotMessages loop through the message channel and send the chatbot message
// to all Clients in the chat room.
func (s *Server) listenChatbotMessages(ctx context.Context) {
	msgCH := make(chan []byte)
	go s.broker.ReadMessage(ctx, msgCH)

	for msg := range msgCH {
		if ctx.Err() == context.Canceled {
			log.Warn("context canceled")
			return
		}

		var output CommandOutput
		if err := json.Unmarshal(msg, &output); err != nil {
			log.WithError(err).Error("could not decode chatbot message")
			continue
		}

		log.WithField("CommandOutput", output).Info("received chatbot message")

		if !s.isValidRoom(output.RoomID) {
			log.Error("chat room not found")
			continue
		}

		msgInput := SendMessageInputEvent{
			Message: output.Message,
			From:    output.From,
			Sent:    time.Now(),
		}

		payload, err := json.Marshal(msgInput)
		if err != nil {
			log.WithError(err).Error("could not encode chatbot message")
			continue
		}

		event := Event{
			Action:  SendMessageAction,
			Payload: payload,
		}

		// broadcast event to all clients in the same chat room
		for client := range s.clients {
			if client.RoomID == output.RoomID {
				client.event <- event
			}
		}
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
		SendMessageAction:        SendMessageHandler,
		JoinRoomAction:           ChatRoomHandler,
		SendChatbotCommandAction: ChatbotCommandHandler,
	}

	return handlers
}
