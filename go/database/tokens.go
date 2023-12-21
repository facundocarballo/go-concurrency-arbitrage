package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/exchange"
	"github.com/facundocarballo/go-concurrency-arbitrage/types/token"
)

func GetAllTokensWithBalance(db *sql.DB, exchanges []exchange.Exchange) []token.Token {
	tokens := GetAllTokens(db)
	tokens = GetAllTokensAmountForEachExchange(tokens, exchanges, db)
	return tokens
}

func GetAllTokens(db *sql.DB) []token.Token {
	rows, err := db.Query(Q_GET_ALL_TOKENS)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tokens []token.Token
	for rows.Next() {
		var token token.Token
		err := rows.Scan(&token.Id, &token.Name, &token.Symbol)
		if err != nil {
			log.Fatal(err)
		}
		token.Amounts = make(map[int]float64)
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tokens
}

func GetAllTokensAmountForEachExchange(
	tokens []token.Token,
	exchanges []exchange.Exchange,
	db *sql.DB,
) []token.Token {
	var newTokens []token.Token
	for _, token := range tokens {
		for _, ex := range exchanges {
			statement, err := db.Prepare(SP_GET_TOKEN_BALANCE_ON_EXCHANGE)
			if err != nil {
				fmt.Printf("Error preparing the stored procedure for get the amount of %s in %s. Error: %s\n", token.Name, ex.Name, err)
				continue
			}
			defer statement.Close()

			var amount float64
			_, err = statement.Exec(ex.Id, token.Id)
			if err != nil {
				fmt.Printf("Error getting the amount of %s in %s. Error: %s\n", token.Symbol, ex.Name, err)
				continue
			}
			err = db.QueryRow("SELECT @amount").Scan(&amount)
			if err != nil {
				log.Fatal(err)
			}
			token.Amounts[ex.Id] = amount
			// fmt.Printf("[%s] %s: %.2f\n", ex.Name, token.Symbol, amount)
		}
		newTokens = append(newTokens, token)
	}
	return newTokens
}
