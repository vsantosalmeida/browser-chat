package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// ErrInvalidRoomID invalid room id
var ErrInvalidRoomID = errors.New("invalid room id")

// Event represents an event to execute some action to the Server or to the Client.
type Event struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler function to execute the required event.
type EventHandler func(event Event, c *Client) error

const (
	// SendMessageAction action to represent a message sent by a Client.
	SendMessageAction = "sendMessage"
	// MessageReceivedAction action to represent a message for a Client read.
	MessageReceivedAction    = "messageReceived"
	JoinRoomAction           = "joinRoom"
	SendChatbotCommandAction = "chatbotCommand"
)

// MessageEvent represents a message sent or received by a user.
type MessageEvent struct {
	Message string    `json:"message"`
	From    string    `json:"from"`
	Sent    time.Time `json:"sent"`
}

// SendMessageHandler handles the client message and send it to all client in the chat room.
//
// if the chat room doesn't exist the event will not be executed.
//
// stores the user message in the DB for the respective chat room.
func SendMessageHandler(event Event, c *Client) error {
	if !c.server.isValidRoom(c.RoomID) {
		return ErrInvalidRoomID
	}

	var input MessageEvent
	if err := json.Unmarshal(event.Payload, &input); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}
	input.Sent = time.Now()

	// error ignored to avoid disconnect a Client
	go c.server.roomUseCase.CreateMessage(c.ID, c.RoomID, input.Message)

	data, err := json.Marshal(input)
	if err != nil {
		return errors.Errorf("could not encode event payload: %v", err)
	}

	output := Event{
		Action:  MessageReceivedAction,
		Payload: data,
	}

	// broadcast event to all clients in the same chat room
	for client := range c.server.clients {
		if client.RoomID == c.RoomID {
			client.event <- output
		}
	}

	return nil
}

type JoinRoomEvent struct {
	RoomID int `json:"roomID"`
}

// ChatRoomHandler if the rooms exist will allow the user to join the chat room.
func ChatRoomHandler(event Event, c *Client) error {
	var joinRoomEvent JoinRoomEvent
	if err := json.Unmarshal(event.Payload, &joinRoomEvent); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}

	if !c.server.isValidRoom(joinRoomEvent.RoomID) {
		return ErrInvalidRoomID
	}

	c.RoomID = joinRoomEvent.RoomID

	log.WithFields(log.Fields{
		"UserID": c.ID,
		"RoomID": c.RoomID,
	}).Info("user joined room")

	return nil
}

// ChatbotCommandEvent command received from a Client.
type ChatbotCommandEvent struct {
	RoomID      int    `json:"roomID"`
	From        string `json:"from"`
	CommandName string `json:"commandName"`
	Command     string `json:"command"`
}

func ChatbotCommandHandler(event Event, c *Client) error {
	var chatbotEvent ChatbotCommandEvent
	// decode the event payload to validate the schema
	if err := json.Unmarshal(event.Payload, &chatbotEvent); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}

	// error ignored to avoid disconnect a Client
	go c.server.broker.WriteMessage(context.Background(), event.Payload)

	return nil
}
