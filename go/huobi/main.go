package huobi

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/facundocarballo/go-concurrency-arbitrage/exchange"
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
)

type Huobi struct {
	ApiKey    string
	SecretKey string
}

func CreateHuobiExchange() *Huobi {
	return &Huobi{
		ApiKey:    "HUOBI_API_KEY",
		SecretKey: "HUOBI_SECRET_KEY",
	}
}

func (exchange *Huobi) GetPrice(pair *exchange.Pair) float64 {

	client := new(client.MarketClient).Init(config.Host)
	resp, err := client.GetLatestTrade(strings.ToLower(pair.Symbol))
	if err != nil {
		fmt.Println("[Huobi] Error al obtener el precio del par:", err)
		os.Exit(1)
	}

	price, err := strconv.ParseFloat(resp.Data[0].Price.String(), 64)
	if err != nil {
		fmt.Println("[Huobi] Error converting the price to float64", err)
		os.Exit(1)
	}

	return price
}

func (exchange *Huobi) GetName() string {
	return "Huobi"
}
