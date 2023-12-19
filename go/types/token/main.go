package token

import "encoding/json"

type Token struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

func CreateToken(id int, name string, symbol string) *Token {
	return &Token{
		Id:     id,
		Name:   name,
		Symbol: symbol,
	}
}

func BodyToToken(body []byte) *Token {
	if len(body) == 0 {
		return nil
	}

	var token Token
	err := json.Unmarshal(body, &token)
	if err != nil {
		return nil
	}

	return &token
}
