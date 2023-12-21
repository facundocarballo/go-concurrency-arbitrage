package core

import (
	"math"
	"testing"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
)

func TestCorrectTime(t *testing.T) {
	timeA := time.Now()
	timeB := timeA.Add(time.Nanosecond * 102)

	resA := CorrectTime(&timeA, &timeB)
	if resA {
		t.Errorf("[CorrectTime] Error. Waiting false, get true. TimeB is 102 ns old than TimeA.")
	}

	resB := CorrectTime(&timeB, &timeA)
	if resB {
		t.Errorf("[CorrectTime] Error. Waiting false, get true. TimeA is 102 ns old than TimeB.")
	}

	timeB = timeA.Add(time.Nanosecond * 50)

	resC := CorrectTime(&timeA, &timeB)
	if !resC {
		t.Errorf("[CorrectTime] Error. Waiting true, get false. TimeB is 50 ns old than TimeA.")
	}

	resD := CorrectTime(&timeB, &timeA)
	if !resD {
		t.Errorf("[CorrectTime] Error. Waiting true, get false. TimeA is 50 ns old than TimeB.")
	}
}

func TestCorrectDifference(t *testing.T) {
	maxPrice := float64(1)
	minPrice := float64(1)

	resA := CorrectDifference(maxPrice, minPrice)

	if resA {
		t.Errorf("[CorrectDifference] Error. Waiting false, get true. Price are the same.")
	}

	maxPrice += 0.00000001

	resB := CorrectDifference(maxPrice, minPrice)

	if resB {
		t.Error("[CorrectDifference] Error. Waiting false, get true. Max price are not %0.01 higher than the min price.")
	}

	resC := CorrectDifference(minPrice, maxPrice)

	if resC {
		t.Errorf("[CorrectDifference] Error. Waiting false, get true. Max price is less than min price.")
	}

	maxPrice += 0.5

	resD := CorrectDifference(maxPrice, minPrice)

	if !resD {
		t.Error("[CorrectDifference] Error. Waiting true, get false. Max price is %0.01 higher than the min price.")
	}
}

func TestGetAmountToBuy(t *testing.T) {
	price := float64(1)
	balance := float64(10)

	resA := GetAmountToBuy(price, balance)
	if resA != balance {
		t.Errorf("[GetAmountToBuy] Error. Waiting 10.00, get %.2f.\n", resA)
	}
}

func TestChangeMaxPrice(t *testing.T) {
	maxPriceScanned := Scanned{
		Price: 0,
	}
	minPriceScanned := Scanned{
		Price: math.MaxFloat64,
	}
	scanned := Scanned{
		Price: 50,
	}

	ChangeMaxPrice(&maxPriceScanned, &minPriceScanned, &scanned)

	if maxPriceScanned.Price != float64(50) {
		t.Errorf("[ChangeMaxPrice] Error. Waiting 50.00 as max price, get %.2f.\n", maxPriceScanned.Price)
	}

	if minPriceScanned.Price != math.MaxFloat64 {
		t.Errorf("[ChangeMaxPrice] Error. Waiting %.2f as min price, get %.2f.\n", math.MaxFloat64, minPriceScanned.Price)
	}

	scanned.Price = float64(101)

	ChangeMaxPrice(&maxPriceScanned, &minPriceScanned, &scanned)

	if maxPriceScanned.Price != float64(101) {
		t.Errorf("[ChangeMaxPrice] Error. Waiting 101.00 as max price, get %.2f.\n", maxPriceScanned.Price)
	}

	if minPriceScanned.Price != float64(50) {
		t.Errorf("[ChangeMaxPrice] Error. Waiting 50.00 as min price, get %.2f.\n", minPriceScanned.Price)
	}
}

func TestGetMaxAndMinPriceScanned(t *testing.T) {
	exchangeA := exchange.Exchange{
		Name: "Binance",
	}
	exchangeB := exchange.Exchange{
		Name: "Bybit",
	}
	maxPriceScanned := Scanned{
		Price: 0,
	}
	minPriceScanned := Scanned{
		Price: math.MaxFloat64,
	}
	scannedA := Scanned{
		Price: 1,
	}
	scannedB := Scanned{
		Price: 1.25,
	}

	mapPair := make(map[exchange.Exchange]Scanned)
	mapPair[exchangeA] = scannedA
	mapPair[exchangeB] = scannedB

	GetMaxAndMinPriceScanned(&maxPriceScanned, &minPriceScanned, mapPair)

	if maxPriceScanned.Price != float64(1.25) {
		t.Errorf("[ChangeMaxPrice] Error. Waiting 1.25 as max price, get %.2f.\n", maxPriceScanned.Price)
	}

	if minPriceScanned.Price != float64(1) {
		t.Errorf("[ChangeMaxPrice] Error. Waiting 1.00 as min price, get %.2f.\n", minPriceScanned.Price)
	}
}
