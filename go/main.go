package main

import (
	"log"
	"sync"

	"github.com/facundocarballo/go-concurrency-arbitrage/scan"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/binance"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange/huobi"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
	"github.com/joho/godotenv"
)

func main() {
	// Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	// Exchanges
	binanceStruct := binance.CreateBinanceExchange()
	huobiStruct := huobi.CreateHuobiExchange()

	exchanges := []exchange.IExchange{binanceStruct, huobiStruct}

	// Pairs
	btcUsdt := pair.CreatePair("BTC", "USDT")
	ethUsdt := pair.CreatePair("ETH", "USDT")
	eosUsdt := pair.CreatePair("EOS", "USDT")
	bnbUsdt := pair.CreatePair("BNB", "USDT")

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		scan.ScanPair(btcUsdt, exchanges)
	}()

	go func() {
		defer wg.Done()
		scan.ScanPair(ethUsdt, exchanges)
	}()

	go func() {
		defer wg.Done()
		scan.ScanPair(eosUsdt, exchanges)
	}()

	go func() {
		defer wg.Done()
		scan.ScanPair(bnbUsdt, exchanges)
	}()

	wg.Wait()
}
