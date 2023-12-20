package database

import (
	"database/sql"
	"log"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/token"
)

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
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tokens
}
