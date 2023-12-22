# Arbitrge Concurrency Bot

## What does this bot?
This bot emulates an arbitrage between multiple exchanges, in this particular case we use:
- Binance
- Huobi
- Bybit
- Bitget
But you can add any exchanges that you want.

The bot gets the prices of a Pair from each exchange, for example, BTC-USDT, and compares the prices.

Getting with the max price and the minimum price.

The idea is that the bot can buy BTC (or other token) in the cheapest exchange and sell it on the more expensive exchange that we compare it.

## How this repo is organized?
📁 db
> Here are all the SQL files. Tables and Stored Procedures that I use to store the data that produces this bot.

📁 go
> Here are all the codes of the bot.

📁 out
> Here is a document that registers all the highest differences in prices between two exchanges in the same pair.

📄 TradeExample.txt
> Explain and simplify how this bot makes trades between two exchanges.