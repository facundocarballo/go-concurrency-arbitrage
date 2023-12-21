package core

import (
	"fmt"
	"os"
	"time"
)

func CorrectTime(timeA *time.Time, timeB *time.Time) bool {
	diff := timeA.Nanosecond() - timeB.Nanosecond()

	return diff < 100 || diff > -100
}

func CorrectDifference(priceA float64, priceB float64) bool {
	return (priceA > priceB*1.0001) || (priceB > priceA*1.0001)
}

func WriteDifferenceFinded(scannedA *Scanned, scannedB *Scanned) {
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

func GetAmountToBuy(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}

func GetAmountToSell(price float64, amount float64) float64 {
	return (amount * PERCENTAGE_OF_USDT) / price
}
