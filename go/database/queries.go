package database

// Stored Procedures
const SP_CREATE_TRADE = "CALL CreateTrade(?, ?, ?, ?, ?)"

// Queries
const Q_GET_ALL_EXCHANGES = "SELECT * FROM Exchange"
const Q_GET_ALL_TOKENS = "SELECT * FROM Token"
