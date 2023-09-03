package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
	readLimit  = 1000
)

type Client struct {
	conn     *websocket.Conn
	server   *Server
	event    chan Event
	ID       int
	Username string
	RoomID   int
}

// NewClient
func NewClient(conn *websocket.Conn, server *Server, username string, id int) *Client {
	return &Client{
		conn:     conn,
		server:   server,
		event:    make(chan Event),
		ID:       id,
		Username: username,
	}
}

// readMessages
func (c *Client) readMessages() {
	defer func() {
		c.server.leave <- c
	}()

	c.conn.SetReadLimit(readLimit)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.handlePong)

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}

		var event Event
		if err = json.Unmarshal(payload, &event); err != nil {
			// decoding errors should not stop the client connection
			continue
		}

		if err = c.server.routeEvent(event, c); err != nil {
			break
		}
	}
}

// writeMessages
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.server.leave <- c
	}()

	for {
		select {
		case event, ok := <-c.event:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
				}
				return
			}

			b, err := json.Marshal(event)
			if err != nil {
				// encoding errors should not stop the client connection
				continue
			}

			if err = c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			}

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}

	}
}

func (c *Client) handlePong(_ string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
