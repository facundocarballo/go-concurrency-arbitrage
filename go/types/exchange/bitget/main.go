package bitget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/facundocarballo/go-concurrency-arbitrage/types/pair"
)

type BitgetResult struct {
	Price string `json:"lastPr"`
}

type BitgetResponse struct {
	Data []BitgetResult `json:"data"`
}

func BodyToBitgetResponse(body []byte) *BitgetResponse {
	if len(body) == 0 {
		return nil
	}

	var res BitgetResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}

	return &res
}

func GetPrice(pair *pair.Pair) float64 {
	url := "https://api.bitget.com/api/v2/spot/market/tickers?symbol=" + pair.GetSymbol()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[Bitget] Error creating the request. %s\n", err)
		return 0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Bybit] Error making the request. %s\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response. %s\n", err)
		return 0
	}

	response := BodyToBitgetResponse(body)
	if response == nil {
		fmt.Printf("Error transforming the response to a data structure.\n")
		return 0
	}

	price, err := strconv.ParseFloat(response.Data[0].Price, 64)
	if err != nil {
		fmt.Printf("Error converting the string to float64.\n")
		return 0
	}

	return price
}
