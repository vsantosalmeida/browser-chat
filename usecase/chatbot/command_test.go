package chatbot_test

import (
	"context"
	"testing"

	"github.com/vsantosalmeida/browser-chat/pkg/stooq"
	"github.com/vsantosalmeida/browser-chat/pkg/stooq/mocks"
	"github.com/vsantosalmeida/browser-chat/usecase/chatbot"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestService_ExecuteCommandStock(t *testing.T) {
	var (
		resp = &stooq.GetStockResponse{
			StockID: "AMZN.US",
			Quote:   "138.12",
		}

		expected = "AMZN.US quote is $138.12 per share"
	)

	stooqAPI := mocks.NewAPI(t)
	svc := chatbot.NewService(nil, stooqAPI, 1)

	stooqAPI.
		On("GetStock", ctx, "amzn.us").
		Return(resp, nil).
		Once()

	result, err := svc.ExecuteCommand(ctx, "stock", "amzn.us")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestService_ExecuteCommandErrors(t *testing.T) {
	var tt = []struct {
		name        string
		commandName string
		command     string
		expected    string
		mockAPIErr  error
	}{
		{
			name:        "When a command handler is not found; should return error",
			commandName: "echo",
			expected:    "invalid command",
		},
		{
			name:        "When stoop API request fail; should return error",
			commandName: "stock",
			command:     "amzn.us",
			expected:    "stooq api error",
			mockAPIErr:  errors.New("stooq api error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			stooqAPI := mocks.NewAPI(t)
			svc := chatbot.NewService(nil, stooqAPI, 1)

			stooqAPI.
				On("GetStock", ctx, "amzn.us").
				Return(nil, tc.mockAPIErr).
				Maybe()

			result, err := svc.ExecuteCommand(ctx, tc.commandName, tc.command)
			assert.EqualError(t, err, tc.expected)
			assert.Empty(t, "", result)
		})
	}
}
