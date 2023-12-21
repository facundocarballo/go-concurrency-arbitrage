package database

// Stored Procedures
const SP_CREATE_TRADE = "CALL CreateTrade(?, ?, ?, ?, ?)"
const SP_GET_TOKEN_BALANCE_ON_EXCHANGE = "CALL GetTokenBalanceOnExchange(?, ?, @amount)"
const SP_GET_TOKEN_BALANCE = "CALL GetTokenBalance(?, ?)"

// Queries
const Q_GET_ALL_EXCHANGES = "SELECT * FROM Exchange"
const Q_GET_ALL_TOKENS = "SELECT * FROM Token"
