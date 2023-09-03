package websocket

import (
	"encoding/json"
	"time"

	"github.com/apex/log"
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

	logger := log.WithFields(log.Fields{
		"UserID": c.ID,
	})

	c.conn.SetReadLimit(readLimit)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.handlePong)

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.WithError(err).Error("unexpected close error")
			}
			return
		}

		var event Event
		if err = json.Unmarshal(payload, &event); err != nil {
			// decoding errors should not stop the client connection
			logger.WithError(err).Error("failed to decode event body")
			continue
		}

		if err = c.server.routeEvent(event, c); err != nil {
			logger.WithError(err).Error("failed to process event")
			return
		}

		logger.WithField("event", event.Action).Info("event sent")
	}
}

// writeMessages
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.server.leave <- c
	}()

	logger := log.WithFields(log.Fields{
		"UserID": c.ID,
	})

	for {
		select {
		case event := <-c.event:
			b, err := json.Marshal(event)
			if err != nil {
				// encoding errors should not stop the client connection
				logger.WithError(err).Error("failed to encode event body")
				continue
			}

			if err = c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				// write errors should not stop the client connection
				logger.WithError(err).Error("failed to write message")
				continue
			}

			logger.Info("event received")

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.WithError(err).Error("ping timeout")
				return
			}

		}

	}
}

func (c *Client) handlePong(_ string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
