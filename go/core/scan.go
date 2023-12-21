package core

import (
	"database/sql"
	"fmt"
	"math"
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

var writeMapMutex sync.Mutex

func GetPrices(
	mapPair map[exchange.Exchange]Scanned,
	pair *pair.Pair,
	exchanges []exchange.Exchange,
	priceGettedChannel *chan bool,
	mapPairMutex *sync.Mutex,
) {
	for {
		mapPairMutex.Lock()
		var wg sync.WaitGroup
		for _, exc := range exchanges {
			wg.Add(1)
			go func(ex exchange.Exchange) {
				defer wg.Done()
				price := ex.GetPrice(pair)
				priceScanned := Scanned{
					Pair:      *pair,
					Exchange:  ex,
					Price:     price,
					Timestamp: time.Now(),
				}
				writeMapMutex.Lock()
				mapPair[ex] = priceScanned
				writeMapMutex.Unlock()
			}(exc)
		}
		wg.Wait()
		*priceGettedChannel <- true
		mapPairMutex.Unlock()
		time.Sleep(time.Second * 15)
	}
}

func AnalizePrice(
	mapPair map[exchange.Exchange]Scanned,
	mapOrders map[Order]bool,
	db *sql.DB,
	priceGettedChannel *chan bool,
	mapPairMutex *sync.Mutex,
) {
	for {
		maxPriceScanned := Scanned{
			Price: 0,
		}
		minPriceScanned := Scanned{
			Price: math.MaxFloat64,
		}

		<-*priceGettedChannel
		mapPairMutex.Lock()
		GetMaxAndMinPriceScanned(&maxPriceScanned, &minPriceScanned, mapPair)
		mapPairMutex.Unlock()

		order := Order{
			ExchangeA: maxPriceScanned.Exchange.Name,
			ExchangeB: minPriceScanned.Exchange.Name,
			Timestamp: time.Now(),
		}
		inverseOrder := Order{
			ExchangeA: minPriceScanned.Exchange.Name,
			ExchangeB: maxPriceScanned.Exchange.Name,
			Timestamp: time.Now(),
		}

		if mapOrders[order] || mapOrders[inverseOrder] {
			fmt.Printf("Order already maked...\n")
			continue
		}

		if MakeOrder(&maxPriceScanned, &minPriceScanned) {
			ExecuteOrder(&minPriceScanned, &maxPriceScanned, db)
			mapOrders[order] = true
		}
	}
}

func ScanPair(pair *pair.Pair, exchanges []exchange.Exchange, db *sql.DB) {
	mapPair := make(map[exchange.Exchange]Scanned)
	mapOrders := make(map[Order]bool)

	var mapPairMutex sync.Mutex
	var priceGettedChannel = make(chan bool)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		GetPrices(mapPair, pair, exchanges, &priceGettedChannel, &mapPairMutex)
	}()

	go func() {
		defer wg.Done()
		AnalizePrice(mapPair, mapOrders, db, &priceGettedChannel, &mapPairMutex)
	}()

	wg.Wait()
}
