-- name: SaveExchangeRates :exec
INSERT INTO exchange_rates (base_currency, rate)
VALUES ($1, $2) ON CONFLICT (base_currency) DO
UPDATE
SET rate = $2;
-- name: GetExchangeRate :one
SELECT *
FROM exchange_rates
WHERE base_currency = $1
LIMIT 1;