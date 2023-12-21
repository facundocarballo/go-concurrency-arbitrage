package core

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type Scanned struct {
	Pair      pair.Pair         `json:"pair"`
	Exchange  exchange.Exchange `json:"exchange"`
	Price     float64           `json:"price"`
	Timestamp time.Time         `json:"timestamp"`
}

var mapPairMutex sync.Mutex

func GetPrices(
	ch chan Scanned,
	mapPair map[exchange.Exchange]Scanned,
	pair *pair.Pair,
	exchanges []exchange.Exchange,
) {
	for {
		for _, exc := range exchanges {
			go func(ex exchange.Exchange) {
				price := ex.GetPrice(pair)
				priceScanned := Scanned{
					Pair:      *pair,
					Exchange:  ex,
					Price:     price,
					Timestamp: time.Now(),
				}
				ch <- priceScanned
				mapPairMutex.Lock()
				mapPair[ex] = priceScanned
				mapPairMutex.Unlock()
			}(exc)
		}
		time.Sleep(time.Second * 15)
	}
}

func AnalizePrice(
	ch chan Scanned,
	mapPair map[exchange.Exchange]Scanned,
	mapOrders map[Order]bool,
	db *sql.DB,
) {
	for {
		priceScanned := <-ch
		mapPairMutex.Lock()
		for _, scanned := range mapPair {
			order := Order{
				ExchangeA: priceScanned.Exchange.Name,
				ExchangeB: scanned.Exchange.Name,
				Timestamp: time.Now(),
			}
			inverseOrder := Order{
				ExchangeA: scanned.Exchange.Name,
				ExchangeB: priceScanned.Exchange.Name,
				Timestamp: time.Now(),
			}

			if mapOrders[order] || mapOrders[inverseOrder] {
				fmt.Printf("Order already maked...\n")
				continue
			}

			if MakeOrder(&priceScanned, &scanned) {
				if priceScanned.Price < scanned.Price {
					ExecuteOrder(&priceScanned, &scanned, db)
				} else {
					ExecuteOrder(&scanned, &priceScanned, db)
				}
				mapOrders[order] = true
			}
		}
		mapPairMutex.Unlock()
	}
}

func ScanPair(pair *pair.Pair, exchanges []exchange.Exchange, db *sql.DB) {
	chPair := make(chan Scanned)
	mapPair := make(map[exchange.Exchange]Scanned)
	mapOrders := make(map[Order]bool)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		GetPrices(chPair, mapPair, pair, exchanges)
	}()

	go func() {
		defer wg.Done()
		AnalizePrice(chPair, mapPair, mapOrders, db)
	}()

	wg.Wait()
}
