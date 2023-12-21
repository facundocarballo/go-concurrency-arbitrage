package core

import (
	"fmt"
	"os"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
)

func CorrectTime(timeA *time.Time, timeB *time.Time) bool {
	diff := timeA.Nanosecond() - timeB.Nanosecond()

	return diff < 100 || diff > -100
}

func CorrectDifference(maxPrice float64, minPrice float64) bool {
	return (maxPrice > minPrice*1.0001)
}

func WriteDifferenceFinded(maxScanned *Scanned, minScanned *Scanned) {
	f, err := os.OpenFile("../out/diff.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error creating the diff.txt file. Error: %s\n", err)
		return
	}

	realDif := (maxScanned.Price / minScanned.Price) - 1

	_, err = fmt.Fprintf(
		f,
		"[%s] DIFFERENCE FINDED. %.2f VS %.2f (%.7f)\n",
		maxScanned.Pair.GetSymbol(),
		maxScanned.Price,
		minScanned.Price,
		realDif,
	)
	if err != nil {
		fmt.Printf("Error writting the diff.txt file. Error: %s\n", err)
		return
	}
	f.Close()
}

func GetAmountToBuy(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}

func GetAmountToSell(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}

func ChangeMaxPrice(
	maxPriceScanned *Scanned,
	minPriceScanned *Scanned,
	scanned *Scanned,
) {
	if maxPriceScanned.Price != 0 && maxPriceScanned.Price < minPriceScanned.Price {
		*minPriceScanned = *maxPriceScanned
	}
	*maxPriceScanned = *scanned
}

func GetMaxAndMinPriceScanned(
	maxPriceScanned *Scanned,
	minPriceScanned *Scanned,
	mapPair map[exchange.Exchange]Scanned,
) {
	for _, scanned := range mapPair {
		if scanned.Price > maxPriceScanned.Price {
			ChangeMaxPrice(maxPriceScanned, minPriceScanned, &scanned)
			continue
		}
		if scanned.Price < minPriceScanned.Price {
			*minPriceScanned = scanned
			continue
		}
	}
}
