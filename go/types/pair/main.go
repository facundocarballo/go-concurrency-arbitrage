package pair

import (
	"github.com/facundocarballo/go-concurrency-arbitrage/types/token"
)

type Pair struct {
	TokenA token.Token `json:"token_a"`
	TokenB token.Token `json:"token_b"`
}

func (pair *Pair) GetSymbol() string {
	return pair.TokenA.Symbol + pair.TokenB.Symbol
}

func CreatePair(tokenA token.Token, tokenB token.Token) *Pair {
	return &Pair{
		TokenA: tokenA,
		TokenB: tokenB,
	}
}

func IsPairCreated(
	tokenA token.Token,
	tokenB token.Token,
	pairMap map[string]bool,
) bool {
	if tokenA.Id == tokenB.Id {
		return true
	}
	return pairMap[tokenB.Symbol+tokenA.Symbol]
}

func GetAllPairs(tokens []token.Token) []Pair {
	var pairs []Pair
	pairMap := make(map[string]bool)
	usdt := tokens[1]
	for _, token := range tokens {
		if IsPairCreated(usdt, token, pairMap) {
			continue
		}
		pair := Pair{
			TokenA: token,
			TokenB: usdt,
		}
		pairs = append(pairs, pair)
		pairMap[token.Symbol+usdt.Symbol] = true
	}

	return pairs
}
