package stooq

import "context"

// API abstracts the request/response to communicate with stooq API.
type API interface {
	GetStock(ctx context.Context, stockID string) (GetStockResponse, error)
}
