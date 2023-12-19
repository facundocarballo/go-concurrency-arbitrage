package exchange

import (
	"encoding/json"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type IExchange interface {
	GetPrice(pair *pair.Pair) float64
	GetName() string
}

type Exchange struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func CreateExchange(id int, name string) *Exchange {
	return &Exchange{
		Id:   id,
		Name: name,
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
