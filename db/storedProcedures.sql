DELIMITER //
CREATE PROCEDURE CreateExchange(IN name VARCHAR(255))
BEGIN
	INSERT INTO Exchange (name) VALUES (name);
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE CreateToken(IN name VARCHAR(255), IN symbol VARCHAR(6))
BEGIN
	INSERT INTO Token (name, symbol) VALUES (name, symbol);
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE CreateTrade(
	IN exchange_id INT, 
	IN token_in_id INT, 
    IN amount_in FLOAT8,
    IN token_out_id INT,  
    IN amount_out FLOAT8
)
BEGIN
	INSERT INTO Trade 
    (exchange_id, token_in, amount_in, token_out, amount_out) 
    VALUES 
    (exchange_id, token_in_id, amount_in, token_out_id, amount_out);
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE GetTokenBalanceOnExchange(
	IN exchange_id INT, 
	IN token_id INT, 
    OUT balance FLOAT8
)
BEGIN
	DECLARE total_amount_in FLOAT8;
    DECLARE total_amount_out FLOAT8;
    
	SELECT 
		COALESCE(SUM(amount_in), 0) AS total_in
	INTO total_amount_in
    FROM Trade
    WHERE Trade.exchange_id = exchange_id
    AND Trade.token_in = token_id;
    
	SELECT 
		COALESCE(SUM(amount_out), 0) AS total_out
	INTO total_amount_out
    FROM Trade
    WHERE Trade.exchange_id = exchange_id
    AND Trade.token_out = token_id;
	
	SET balance = total_amount_in - total_amount_out;
    
    IF total_amount_in IS NULL AND total_amount_out IS NULL THEN
        SET balance = 0;
    END IF;
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE GetTokenBalance(
	IN token_id INT, 
    OUT balance FLOAT8
)
BEGIN
	DECLARE total_amount_in FLOAT8;
    DECLARE total_amount_out FLOAT8;
    
	SELECT 
		COALESCE(SUM(amount_in), 0) AS total_in
	INTO total_amount_in
    FROM Trade
    WHERE Trade.token_in = token_id;
    
	SELECT 
		COALESCE(SUM(amount_out), 0) AS total_out
	INTO total_amount_out
    FROM Trade
    WHERE Trade.token_out = token_id;
	
	SET balance = total_amount_in - total_amount_out;
    
    IF total_amount_in IS NULL AND total_amount_out IS NULL THEN
        SET balance = 0;
    END IF;
END //
DELIMITER ;