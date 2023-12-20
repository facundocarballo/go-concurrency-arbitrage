package huobi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
)

func GetPrice(pair *pair.Pair) float64 {
	client := new(client.MarketClient).Init(config.Host)
	resp, err := client.GetLatestTrade(strings.ToLower(pair.GetSymbol()))
	if err != nil {
		fmt.Printf("[Huobi] Cannot get the price of this pair: %s. Error: %s\n", pair.GetSymbol(), err)
		return 0
	}

	price, err := strconv.ParseFloat(resp.Data[0].Price.String(), 64)
	if err != nil {
		fmt.Printf("[Huobi] Error converting the price to float64. %s\n", err)
		return 0
	}

	return price
}
