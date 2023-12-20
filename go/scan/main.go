package scan

import (
	"fmt"
	"sync"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type Scaned struct {
	Pair      pair.Pair          `json:"pair"`
	Exchange  exchange.IExchange `json:"exchange"`
	Price     float64            `json:"price"`
	Timestamp time.Time          `json:"timestamp"`
}

type Order struct {
	ExchangeA string    `json:"exchange_a"`
	ExchangeB string    `json:"exchange_b"`
	Timestamp time.Time `json:"timestamp"`
}

func CorrectTime(timeA *time.Time, timeB *time.Time) bool {
	diff := timeA.Nanosecond() - timeB.Nanosecond()

	return diff < 100 || diff > -100
}

func MakeOrder(scannedA *Scaned, scannedB *Scaned) bool {
	if scannedA.Price == 0 || scannedB.Price == 0 {
		return false
	}

	if scannedA.Price == scannedB.Price {
		return false
	}

	return CorrectTime(&scannedA.Timestamp, &scannedB.Timestamp)
}

func GetPrices(
	ch chan Scaned,
	mapPair map[exchange.IExchange]Scaned,
	pair *pair.Pair,
	exchanges []exchange.IExchange,
) {
	for {
		for _, exc := range exchanges {
			go func(ex exchange.IExchange) {
				price := ex.GetPrice(pair)
				priceScanned := Scaned{
					Pair:      *pair,
					Exchange:  ex,
					Price:     price,
					Timestamp: time.Now(),
				}
				ch <- priceScanned
				mapPair[ex] = priceScanned
			}(exc)
		}
		time.Sleep(time.Second * 15)
	}
}

func AnalizePrice(
	ch chan Scaned,
	mapPair map[exchange.IExchange]Scaned,
	mapOrders map[Order]bool,
) {
	for {
		priceScanned := <-ch

		for ex, scanned := range mapPair {
			order := Order{
				ExchangeA: priceScanned.Exchange.GetName(),
				ExchangeB: scanned.Exchange.GetName(),
				Timestamp: time.Now(),
			}
			inverseOrder := Order{
				ExchangeA: scanned.Exchange.GetName(),
				ExchangeB: priceScanned.Exchange.GetName(),
				Timestamp: time.Now(),
			}

			if mapOrders[order] || mapOrders[inverseOrder] {
				fmt.Printf("Order already maked...\n")
				continue
			}

			if MakeOrder(&priceScanned, &scanned) {
				// TODO:
				// Ejecutar la orden de un exchange.
				// Llevar el resultado al otro exchange.
				// Ejectuar la orden del otro exchange.
				// Enviar el resultado positivo al exchange original.
				fmt.Printf(
					"(%s) in [%s] is $%f and in [%s] is $%f\n",
					priceScanned.Pair.GetSymbol(),
					priceScanned.Exchange.GetName(),
					priceScanned.Price,
					ex.GetName(),
					scanned.Price,
				)
				mapOrders[order] = true
			}
		}
	}
}

func ScanPair(pair *pair.Pair, exchanges []exchange.IExchange) {
	chPair := make(chan Scaned)
	mapPair := make(map[exchange.IExchange]Scaned)
	mapOrders := make(map[Order]bool)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		GetPrices(chPair, mapPair, pair, exchanges)
	}()

	go func() {
		defer wg.Done()
		AnalizePrice(chPair, mapPair, mapOrders)
	}()

	wg.Wait()
}
