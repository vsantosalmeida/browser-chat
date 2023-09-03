package chatbot

import (
	"context"

	"github.com/vsantosalmeida/browser-chat/pkg/stooq"
)

// Service implements UseCase interface.
type Service struct {
	broker   Broker
	stock    stooq.API
	workers  int
	handlers map[string]CommandHandler
}

// NewService Service builder.
func NewService(broker Broker, stock stooq.API, workers int) *Service {
	// avoids a zero or negative workers number
	if workers <= 0 {
		workers = 1
	}

	svc := &Service{
		broker:  broker,
		stock:   stock,
		workers: workers,
	}
	svc.handlers = svc.initHandlers()

	return svc
}

// Start starts the chatbot workers to process commands.
func (s *Service) Start(ctx context.Context) {
	for i := 0; i < s.workers; i++ {
		w := &Worker{
			id:  i,
			svc: s,
		}

		go w.StartAndConsume(ctx)
	}
}
