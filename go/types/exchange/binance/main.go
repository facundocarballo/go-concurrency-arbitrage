package binance

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type Binance struct {
	ApiKey    string
	SecretKey string
}

func CreateBinanceExchange() *Binance {
	return &Binance{
		ApiKey:    "BINANCE_API_KEY",
		SecretKey: "BINANCE_SECRET_KEY",
	}
}

func (exchange *Binance) GetPrice(pair *pair.Pair) float64 {

	apiKey := os.Getenv(exchange.ApiKey)
	secretKey := os.Getenv(exchange.SecretKey)

	client := binance.NewClient(apiKey, secretKey)

	ticker, err := client.NewListPricesService().Symbol(pair.Symbol).Do(context.Background())
	if err != nil {
		fmt.Println("[Binance] Error al obtener el precio del par:", err)
		os.Exit(1)
	}

	price, err := strconv.ParseFloat(ticker[0].Price, 64)
	if err != nil {
		fmt.Println("[Binance] Error converting the price to float64", err)
		os.Exit(1)
	}

	return price
}

func (exchange *Binance) GetName() string {
	return "Binance"
}
