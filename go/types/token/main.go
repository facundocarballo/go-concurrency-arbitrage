package token

import "encoding/json"

type Token struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Symbol  string          `json:"symbol"`
	Amounts map[int]float64 `json:"amounts"`
}

func CreateToken(
	id int,
	name string,
	symbol string,
	amounts map[int]float64,
) *Token {
	return &Token{
		Id:      id,
		Name:    name,
		Symbol:  symbol,
		Amounts: amounts,
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

func (token *Token) GetAmountOnExchange(exchangeId int) float64 {
	return token.Amounts[exchangeId]
}
