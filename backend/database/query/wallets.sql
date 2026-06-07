-- Manage wallets
-- name: CreateWallet :one
INSERT INTO wallets (user_id, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWalletByID :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletByIDForUpdate :one
-- FOR NO KEY UPDATE will lock the selected row for update, but allow other transactions to read it. This is useful when we want to update the wallet balance after checking the current balance.
SELECT * FROM wallets
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: AddWalletBalance :one
UPDATE wallets
SET balance = balance + $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- Transfer money between 2 wallets
-- name: CreateTransfer :one
INSERT INTO transfers (from_wallet_id, to_wallet_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetListTransfers :many
SELECT * FROM transfers
WHERE from_wallet_id = $1 OR to_wallet_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- Record entries for auditing
-- name: CreateEntry :one
INSERT INTO entries (wallet_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetListEntries :many
SELECT * FROM entries
WHERE wallet_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: DeleteWalletByID :exec
DELETE FROM wallets
WHERE id = $1 AND user_id = $2;