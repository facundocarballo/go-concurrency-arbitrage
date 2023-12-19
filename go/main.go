package main

import (
	"log"
	"sync"

	"github.com/facundocarballo/go-concurrency-arbitrage/binance"
	"github.com/facundocarballo/go-concurrency-arbitrage/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/huobi"
	"github.com/facundocarballo/go-concurrency-arbitrage/scan"
	"github.com/joho/godotenv"
)

func main() {
	// Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Exchanges
	binanceStruct := binance.CreateBinanceExchange()
	huobiStruct := huobi.CreateHuobiExchange()

	exchanges := []exchange.IExchange{binanceStruct, huobiStruct}

	// Pairs
	btcUsdt := exchange.CreatePair("BTC", "USDT")
	ethUsdt := exchange.CreatePair("ETH", "USDT")
	eosUsdt := exchange.CreatePair("EOS", "USDT")
	bnbUsdt := exchange.CreatePair("BNB", "USDT")

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
