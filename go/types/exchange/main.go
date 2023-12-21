package exchange

import (
	"encoding/json"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/binance"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/bitget"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/bybit"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/huobi"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type Exchange struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

func CreateExchange(id int, name string, apiKey string, secretKey string) *Exchange {
	return &Exchange{
		Id:        id,
		Name:      name,
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}
}

func BodyToExchange(body []byte) *Exchange {
	if len(body) == 0 {
		return nil
	}

	var exchange Exchange
	err := json.Unmarshal(body, &exchange)
	if err != nil {
		return nil
	}

	return &exchange
}

func (exchange *Exchange) GetPrice(pair *pair.Pair) float64 {
	switch exchange.Id {
	case 1:
		return binance.GetPrice(pair, exchange.ApiKey, exchange.SecretKey)
	case 2:
		return huobi.GetPrice(pair)
	case 4:
		return bybit.GetPrice(pair)
	case 5:
		return bitget.GetPrice(pair)
	}
	return 0
}
