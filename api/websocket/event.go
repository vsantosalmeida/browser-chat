package websocket

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Event struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler
type EventHandler func(event Event, c *Client) error

const (
	SendMessageAction = "sendMessage"
	JoinRoomAction    = "joinRoom"
)

// SendMessageInputEvent
type SendMessageInputEvent struct {
	Message string    `json:"message"`
	From    string    `json:"from"`
	Sent    time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *Client) error {
	if !c.server.isValidRoom(c.RoomID) {
		return errors.New("invalid room id")
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

func ChatRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return errors.Errorf("could not decode event payload: %v", err)
	}

	if !c.server.isValidRoom(changeRoomEvent.RoomID) {
		return errors.New("invalid room id")
	}

	c.RoomID = changeRoomEvent.RoomID

	return nil
}
