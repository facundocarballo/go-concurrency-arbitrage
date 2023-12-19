package exchange

type IExchange interface {
	GetPrice(pair *Pair) float64
	GetName() string
}

type Pair struct {
	TokenA string `json:"token_a"`
	TokenB string `json:"token_b"`
	Symbol string `json:"symbol"`
	// Exchange IExchange `json:"exchange"`
}

func CreatePair(tokenA string, tokenB string) *Pair {
	return &Pair{
		TokenA: tokenA,
		TokenB: tokenB,
		Symbol: tokenA + tokenB,
		// Exchange: exchange,
	}
}

// func (pair *Pair) GetPrice() string {
// 	return pair.Exchange.GetPrice(pair)
// }
