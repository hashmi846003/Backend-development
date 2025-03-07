-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1
For UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: UpdateAccount :exec
UPDATE accounts
SET balance=$2
WHERE id=$1
RETURNING id,owner,balance,currency,created_at;
-- name: DeleteAccount :exec
DELETE FROM accounts
where id=$1;