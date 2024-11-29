package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"garantex/internal/domain/entity"
	"net/http"
)

type GarantexClient struct {
	BaseURL string
}

func New(baseURL string) *GarantexClient {
	return &GarantexClient{
		BaseURL: baseURL,
	}
}

func (c *GarantexClient) GetTrades() ([]entity.Trade, error) {

	url := fmt.Sprintf("%s/api/v2/trades?market=%s", c.BaseURL, "usdtrub")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch trades: status %d", resp.StatusCode)
	}

	var trades []entity.Trade
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return nil, errors.New("no trades found")
	}

	return trades, nil
}
