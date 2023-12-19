CREATE DATABASE GO_CONCURRENCY_ARBITRAGE;
USE GO_CONCURRENCY_ARBITRAGE;

CREATE TABLE Exchange (
	id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE Token (
	id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    symbol VARCHAR(6)
);

CREATE TABLE Trade (
	id INT AUTO_INCREMENT PRIMARY KEY,
    exchange_id INT NOT NULL,
    token_in INT,
    amount_in FLOAT8,
    token_out INT,
    amount_out FLOAT8,
    FOREIGN KEY (token_in) REFERENCES Token(id),
    FOREIGN KEY (token_out) REFERENCES Token(id),
    FOREIGN KEY (exchange_id) REFERENCES Exchange(id)
);