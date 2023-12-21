package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/facundocarballo/go-concurrency-arbitrage/core"
	"github.com/facundocarballo/go-concurrency-arbitrage/database"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
	"github.com/joho/godotenv"
)

func main() {
	// Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	db, err := sql.Open("mysql", database.GetDSN())
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	exchanges := database.GetAllExchanges(db)
	tokens := database.GetAllTokensWithBalance(db, exchanges)
	pairs := pair.GetAllPairs(tokens)

	var wg sync.WaitGroup
	for _, p := range pairs {
		wg.Add(1)
		go func(p pair.Pair) {
			defer wg.Done()
			core.ScanPair(&p, exchanges, db)
		}(p)
	}
	wg.Wait()
}
