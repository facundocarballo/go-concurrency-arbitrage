package trade

import "encoding/json"

type Trade struct {
	Id         int     `json:"id"`
	ExchangeId int     `json:"exchange_id"`
	TokenIn    int     `json:"token_in"`
	AmountIn   float64 `json:"amount_in"`
	TokenOut   int     `json:"token_out"`
	AmountOut  float64 `json:"amount_out"`
}

func CreateTrade(
	id int,
	exchangeId int,
	tokenIn int,
	amountIn float64,
	tokenOut int,
	amountOut float64,
) *Trade {
	return &Trade{
		Id:         id,
		ExchangeId: exchangeId,
		TokenIn:    tokenIn,
		AmountIn:   amountIn,
		TokenOut:   tokenOut,
		AmountOut:  amountOut,
	}
}

func BodyToTrade(body []byte) *Trade {
	if len(body) == 0 {
		return nil
	}

	var trade Trade
	err := json.Unmarshal(body, &trade)
	if err != nil {
		return nil
	}

	return &trade
}

func ExecuteTrade() {

}
