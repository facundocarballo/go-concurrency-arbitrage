package pair

type Pair struct {
	TokenA string `json:"token_a"`
	TokenB string `json:"token_b"`
	Symbol string `json:"symbol"`
}

func CreatePair(tokenA string, tokenB string) *Pair {
	return &Pair{
		TokenA: tokenA,
		TokenB: tokenB,
		Symbol: tokenA + tokenB,
	}
}
