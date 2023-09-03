package stooq

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	baseURL    = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	stockIdIdx = 0
	quoteIdx   = 6
)

// Client implements API.
type Client struct {
	httpClient *http.Client
}

type GetStockResponse struct {
	StockID string
	Quote   string
}

// NewClient Client builder.
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// GetStock retrieves the stock quote close value.
func (c *Client) GetStock(ctx context.Context, stockID string) (*GetStockResponse, error) {
	url := fmt.Sprintf(baseURL, stockID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("expected 200 status code, got %d", resp.StatusCode)
	}

	return buildResponse(resp.Body)
}

// buildResponse read the csv response from stooq API and parses to a GetStockResponse.
func buildResponse(r io.ReadCloser) (*GetStockResponse, error) {
	data, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read csv response")
	}

	resp := &GetStockResponse{}

	for i, line := range data {
		if i > 0 { // omit header line
			for j, field := range line {
				if j == stockIdIdx {
					resp.StockID = field
				} else if j == quoteIdx {
					resp.Quote = field
				}
			}
		}
	}
	return resp, nil
}
