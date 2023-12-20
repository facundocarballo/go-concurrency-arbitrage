package binance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

func GetPrice(pair *pair.Pair, apiKey string, secretKey string) float64 {
	client := binance.NewClient(apiKey, secretKey)

	ticker, err := client.NewListPricesService().Symbol(pair.GetSymbol()).Do(context.Background())
	if err != nil {
		fmt.Printf("[Binance] Cannot get the price of this pair: %s. Error: %s\n", pair.GetSymbol(), err)
		return 0
	}

	price, err := strconv.ParseFloat(ticker[0].Price, 64)
	if err != nil {
		fmt.Println("[Binance] Error converting the price to float64", err)
		return 0
	}

	return price
}
