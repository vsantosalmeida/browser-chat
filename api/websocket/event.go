package websocket

import (
	"encoding/json"
	"time"

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
	SendMessageAction = "sendMessage"
	JoinRoomAction    = "joinRoom"
)

// SendMessageInputEvent represents a message sent by a user to the chat room.
type SendMessageInputEvent struct {
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

	var input SendMessageInputEvent
	if err := json.Unmarshal(event.Payload, &input); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}
	input.Sent = time.Now()

	go c.server.roomUseCase.CreateMessage(c.ID, c.RoomID, input.Message)

	data, err := json.Marshal(input)
	if err != nil {
		return errors.Errorf("could not encode event payload: %v", err)
	}

	var output Event
	output.Payload = data
	output.Action = SendMessageAction

	for client := range c.server.clients {
		if client.RoomID == c.RoomID {
			client.event <- output
		}
	}

	return nil
}

type ChangeRoomEvent struct {
	RoomID int `json:"roomID"`
}

// ChatRoomHandler if the rooms exist will allow the user to join the chat room.
func ChatRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}

	if !c.server.isValidRoom(changeRoomEvent.RoomID) {
		return ErrInvalidRoomID
	}

	c.RoomID = changeRoomEvent.RoomID

	return nil
}
