package database

import (
	"database/sql"
	"log"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	_ "github.com/go-sql-driver/mysql"
)

func GetAllExchanges(db *sql.DB) []exchange.Exchange {
	rows, err := db.Query(Q_GET_ALL_EXCHANGES)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exchanges []exchange.Exchange
	for rows.Next() {
		var exchange exchange.Exchange
		err := rows.Scan(&exchange.Id, &exchange.Name, &exchange.ApiKey, &exchange.SecretKey)
		if err != nil {
			log.Fatal(err)
		}
		exchanges = append(exchanges, exchange)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return exchanges
}
