package core

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/trade"
)

type Order struct {
	ExchangeA string    `json:"exchange_a"`
	ExchangeB string    `json:"exchange_b"`
	Timestamp time.Time `json:"timestamp"`
}

const PERCENTAGE_OF_USDT = 1

var mutex sync.Mutex

func MakeOrder(maxScanned *Scanned, minScanned *Scanned) bool {
	if maxScanned.Price == 0 || minScanned.Price == 0 {
		return false
	}

	if !CorrectDifference(maxScanned.Price, minScanned.Price) {
		return false
	} else {
		WriteDifferenceFinded(maxScanned, minScanned)
	}

	return CorrectTime(&maxScanned.Timestamp, &minScanned.Timestamp)
}

func BuyTokenCheapExchange(
	cheapScan *Scanned,
	expensiveScan *Scanned,
	db *sql.DB,
) trade.Trade {
	transaction := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    cheapScan.Pair.TokenA.Id,
		AmountIn:   GetAmountToBuy(cheapScan.Price, cheapScan.Pair.Usdt.GetAmountOnExchange(cheapScan.Exchange.Id)),
		TokenOut:   cheapScan.Pair.Usdt.Id,
		AmountOut:  cheapScan.Pair.Usdt.GetAmountOnExchange(cheapScan.Exchange.Id) * PERCENTAGE_OF_USDT,
	}
	transaction.Execute(db)
	fmt.Printf(
		"[%s] Bought %f %s at $%f\n",
		cheapScan.Exchange.Name,
		transaction.AmountIn,
		cheapScan.Pair.TokenA.Symbol,
		cheapScan.Price,
	)
	return transaction
}

func TransferTokenBoughtToExpensiveExchange(
	cheapScan *Scanned,
	expensiveScan *Scanned,
	db *sql.DB,
	previousTrade trade.Trade,
) trade.Trade {
	// Reduce the tokens boughted in the cheap exchange.
	transaction := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    cheapScan.Pair.TokenA.Id,
		AmountIn:   0,
		TokenOut:   cheapScan.Pair.TokenA.Id,
		AmountOut:  previousTrade.AmountIn,
	}
	transaction.Execute(db)

	// Increment the tokens boughted in the expensive exchange.
	secondTransaction := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.TokenA.Id,
		AmountIn:   transaction.AmountOut,
		TokenOut:   cheapScan.Pair.TokenA.Id,
		AmountOut:  0,
	}
	secondTransaction.Execute(db)

	return secondTransaction
}

func SellTokenInExpensiveExchange(
	cheapScan *Scanned,
	expensiveScan *Scanned,
	db *sql.DB,
	previousTrade trade.Trade,
) trade.Trade {
	transaction := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   expensiveScan.Price * previousTrade.AmountIn,
		TokenOut:   expensiveScan.Pair.TokenA.Id,
		AmountOut:  previousTrade.AmountIn,
	}
	transaction.Execute(db)
	fmt.Printf(
		"[%s] Sold %f %s at $%f\n",
		expensiveScan.Exchange.Name,
		transaction.AmountOut,
		expensiveScan.Pair.TokenA.Symbol,
		expensiveScan.Price,
	)
	return transaction
}

func TransferTheProfiToTheCheapExchange(
	cheapScan *Scanned,
	expensiveScan *Scanned,
	db *sql.DB,
	previousTrade trade.Trade,
) {
	// Remove the USDT earned on the previous sell, in the expensive Exchange
	transaction := trade.Trade{
		Id:         1,
		ExchangeId: expensiveScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   0,
		TokenOut:   expensiveScan.Pair.Usdt.Id,
		AmountOut:  previousTrade.AmountIn,
	}
	transaction.Execute(db)

	// Increment the USDT earned on the previous sell, in the cheap Exchange
	secondTransaction := trade.Trade{
		Id:         1,
		ExchangeId: cheapScan.Exchange.Id,
		TokenIn:    expensiveScan.Pair.Usdt.Id,
		AmountIn:   transaction.AmountOut,
		TokenOut:   expensiveScan.Pair.Usdt.Id,
		AmountOut:  0,
	}
	secondTransaction.Execute(db)
}

func ExecuteOrder(
	cheapScan *Scanned,
	expensiveScan *Scanned,
	db *sql.DB,
) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Printf("--------- Begin Transaction %s ---------\n", cheapScan.Pair.GetSymbol())

	firstTransaction := BuyTokenCheapExchange(cheapScan, expensiveScan, db)
	secondTransaction := TransferTokenBoughtToExpensiveExchange(cheapScan, expensiveScan, db, firstTransaction)
	thirdTransaction := SellTokenInExpensiveExchange(cheapScan, expensiveScan, db, secondTransaction)
	TransferTheProfiToTheCheapExchange(cheapScan, expensiveScan, db, thirdTransaction)

	fmt.Printf("Profit of $%.7f\n", thirdTransaction.AmountIn-firstTransaction.AmountOut)
	fmt.Printf("---------- End Transaction %s ----------\n", cheapScan.Pair.GetSymbol())
}
