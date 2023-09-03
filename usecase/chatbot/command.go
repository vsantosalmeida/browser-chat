package chatbot

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// ErrInvalidCommand invalid command
var ErrInvalidCommand = errors.New("invalid command")

const stockCommand = "stock"

// CommandHandler function to receive and execute a command.
type CommandHandler func(ctx context.Context, command string) (string, error)

// CommandInput command received from a user.
type CommandInput struct {
	RoomID      int    `json:"roomID"`
	From        string `json:"from"`
	CommandName string `json:"commandName"`
	Command     string `json:"command"`
}

// CommandOutput result of executed command.
type CommandOutput struct {
	RoomID  int    `json:"roomID"`
	From    string `json:"from"`
	Message string `json:"message"`
}

// ExecuteCommand finds the CommandHandler and execute the command.
// if a handler isn't found returns an error.
func (s *Service) ExecuteCommand(ctx context.Context, commandName, command string) (string, error) {
	if h, ok := s.handlers[commandName]; ok {
		return h(ctx, command)
	}
	return "", ErrInvalidCommand
}

// stockCommandHandler retrieves the stock quotes for a given stock ID.
func (s *Service) stockCommandHandler(ctx context.Context, command string) (string, error) {
	resp, err := s.stock.GetStock(ctx, command)
	if err != nil {
		log.WithFields(log.Fields{
			"command": command,
			"error":   err,
		}).Error("failed to retrieve stock quote")
		return "", err
	}

	return fmt.Sprintf("%s quote is $%s per share", resp.StockID, resp.Quote), nil
}

func (s *Service) initHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		stockCommand: s.stockCommandHandler,
	}
}
