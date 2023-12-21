package scan

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/trade"
)

type Scaned struct {
	Pair      pair.Pair         `json:"pair"`
	Exchange  exchange.Exchange `json:"exchange"`
	Price     float64           `json:"price"`
	Timestamp time.Time         `json:"timestamp"`
}

type Order struct {
	ExchangeA string    `json:"exchange_a"`
	ExchangeB string    `json:"exchange_b"`
	Timestamp time.Time `json:"timestamp"`
}

const PERCENTAGE_OF_USDT = 1

var orderMutex sync.Mutex
var mapPairMutex sync.Mutex

func CorrectTime(timeA *time.Time, timeB *time.Time) bool {
	diff := timeA.Nanosecond() - timeB.Nanosecond()

	return diff < 100 || diff > -100
}

func CorrectDifference(priceA float64, priceB float64) bool {
	return (priceA > priceB*1.0001) || (priceB > priceA*1.0001)
}

func WriteDifferenceFinded(scannedA *Scaned, scannedB *Scaned) {
	f, err := os.OpenFile("../out/diff.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error creating the diff.txt file. Error: %s\n", err)
		return
	}

	var realDif float64
	if scannedA.Price > scannedB.Price {
		realDif = (scannedA.Price / scannedB.Price) - 1
	} else {
		realDif = (scannedB.Price / scannedA.Price) - 1
	}

	_, err = fmt.Fprintf(f, "[%s] DIFFERENCE FINDED. %.2f VS %.2f (%.7f)\n", scannedA.Pair.GetSymbol(), scannedA.Price, scannedB.Price, realDif)
	if err != nil {
		fmt.Printf("Error writting the diff.txt file. Error: %s\n", err)
		return
	}
	f.Close()
}

func MakeOrder(scannedA *Scaned, scannedB *Scaned) bool {
	if scannedA.Price == 0 || scannedB.Price == 0 {
		return false
	}

	if scannedA.Price == scannedB.Price {
		return false
	}

	if !CorrectDifference(scannedA.Price, scannedB.Price) {
		return false
	} else {
		WriteDifferenceFinded(scannedA, scannedB)
	}

	return CorrectTime(&scannedA.Timestamp, &scannedB.Timestamp)
}

func GetPrices(
	ch chan Scaned,
	mapPair map[exchange.Exchange]Scaned,
	pair *pair.Pair,
	exchanges []exchange.Exchange,
) {
	for {
		for _, exc := range exchanges {
			go func(ex exchange.Exchange) {
				price := ex.GetPrice(pair)
				priceScanned := Scaned{
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
	ch chan Scaned,
	mapPair map[exchange.Exchange]Scaned,
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

func GetAmountToBuy(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}

func GetAmountToSell(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}

func ExecuteOrder(cheapScan *Scaned, expensiveScan *Scaned, db *sql.DB) {
	orderMutex.Lock()
	defer orderMutex.Unlock()
	fmt.Printf("--------- Begin Transaction %s ---------\n", cheapScan.Pair.GetSymbol())
	// Buy in the cheap exchange the token.
	firstStep := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    cheapScan.Pair.TokenA.Id,
		AmountIn:   GetAmountToBuy(cheapScan.Price, cheapScan.Pair.Usdt.GetAmountOnExchange(cheapScan.Exchange.Id)),
		TokenOut:   cheapScan.Pair.Usdt.Id,
		AmountOut:  cheapScan.Pair.Usdt.GetAmountOnExchange(cheapScan.Exchange.Id) * PERCENTAGE_OF_USDT,
	}
	firstStep.Execute(db)
	fmt.Printf("[%s] Bought %f %s at $%f\n", cheapScan.Exchange.Name, firstStep.AmountIn, cheapScan.Pair.TokenA.Symbol, cheapScan.Price)

	// Reduce the tokens boughted in the cheap exchange.
	secondStep := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    cheapScan.Pair.TokenA.Id,
		AmountIn:   0,
		TokenOut:   cheapScan.Pair.TokenA.Id,
		AmountOut:  firstStep.AmountIn,
	}
	secondStep.Execute(db)

	// Increment the tokens boughted in the expensive exchange.
	thirdStep := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.TokenA.Id,
		AmountIn:   secondStep.AmountOut,
		TokenOut:   cheapScan.Pair.TokenA.Id,
		AmountOut:  0,
	}
	thirdStep.Execute(db)

	// Sell the tokens boughted in the expensive exchange.
	fourStep := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   expensiveScan.Price * thirdStep.AmountIn,
		TokenOut:   expensiveScan.Pair.TokenA.Id,
		AmountOut:  thirdStep.AmountIn,
	}
	fourStep.Execute(db)
	fmt.Printf("[%s] Sold %f %s at $%f\n", expensiveScan.Exchange.Name, fourStep.AmountOut, expensiveScan.Pair.TokenA.Symbol, expensiveScan.Price)

	// Remove the USDT earned on the previous sell, in the expensive Exchange
	fiveStep := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   0,
		TokenOut:   expensiveScan.Pair.Usdt.Id,
		AmountOut:  fourStep.AmountIn,
	}
	fiveStep.Execute(db)

	// Increment the USDT earned on the previous sell, in the cheap Exchange
	sixStep := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   fiveStep.AmountOut,
		TokenOut:   expensiveScan.Pair.Usdt.Id,
		AmountOut:  0,
	}
	sixStep.Execute(db)

	fmt.Printf("Profit of $%.7f\n", fourStep.AmountIn-firstStep.AmountOut)
	fmt.Printf("---------- End Transaction %s ----------\n", cheapScan.Pair.GetSymbol())
}

func ScanPair(pair *pair.Pair, exchanges []exchange.Exchange, db *sql.DB) {
	chPair := make(chan Scaned)
	mapPair := make(map[exchange.Exchange]Scaned)
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
