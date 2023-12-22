# Arbitrge Concurrency Bot

## What does this bot?
This bot emulates an arbitrage between multiples exchanges, in this particular case we use:
- Binance
- Huobi
- Bybit
- Bitget
But you can add any exchanges that you want.

The bot get the prices of a Pair from each exchange, for example BTC-USDT and compare the prices.
Getting with the max price and the minimun price.

The idea is that the bot can buy BTC (or other token) in the cheapest exchange and sell it on the more expensive exchange that we compare it.

## How this repo is organized?
📁 db
> Here are all the sql files. Tables and Stored Procedures that I use to store the data that produces this bot.

📁 go
> Here are all the code of the bot.

📁 out
> Here are a document that register all the highest differences of prices between two exchanges in the same pair.

📄 TradeExample.txt
> Explain simplify how this bot make trades between two exchanges.