# Golang Concurrency Arbitrage Bot

## How Works?
Using 7 tokens
- BTC
- ETH
- BNB
- SOL
- XRP
- MATIC

The Bot analyzes the prices of each token in USDT in 4 different exchanges
- Binance
- Huobi
- Bybit
- Bitget

The Bot buys a token in the cheap exchange and sells it on the more expensive exchange.

## Why is a Concurrent Bot?
This is a concurrent Bot because, for each pair that this Bot analyzes, it executes a goroutine that will handle this particular pair of tokens.

Then, in each of those goroutines, for each exchange that this Bot handles execute another goroutine that will be in charge of getting the price for this particular pair of tokens in that exchange.

When a pair of tokens get all the prices of all the exchanges that this Bot handles, another goroutine takes place and will try to analyze the prices getted

Searching for some % difference between the max price of that pair of tokens in a particular exchange and the minimum price of the same pair in another exchange.


## How Can I replicate this Bot?
1. Clone this repo in your computer
```bash
    git clone https://github.com/facundocarballo/go-concurrency-arbitrage.git
```

2. Create a MySQL Database
> Using the same tables and stored procedures that you will find in the '../db' folder.

3. Create Accounts in the exchanges that you will use
> In this case we use
- Binance
- Huobi
- Bybit
- Bitget

4. Create an API Key and Secret Key for each exchange that you want to use.

5. Create an instance of each exchange in your database.
> Example for Binance
API Key: "Hello Everyone 👋🏼"
Secret Key: "We love go ❤️"


```sql
    CALL CreateExchange("Binance", "Hello Everyone 👋🏼", "We love go ❤️");
```

6. Create an instance of each token in your database.
> Example for BTC and USDT

```sql
    CALL CreateToken("Bitcoin", "BTC")
```
```sql
    CALL CreateToken("USD Tether", "USDT")
```

7. Create an **.env** file with this data
```.env
    DB_HOST="localhost"
    DB_PORT="3306"
    DB_PASSWORD="Your Password"
    DB_USER="root"
    DB_NAME="GO_CONCURRENCY_ARBITRAGE"
```
> You can find this data on your MySQL workbench, or just ask ChatGPT how you can get it 🤣

8. Install all the dependencies of this Golang Project
```bash
    go get -u ./...
```

9. Run the Bot
```bash
    cd go && go run main.go
```