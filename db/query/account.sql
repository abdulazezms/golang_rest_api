-- name: CreateAccount :one
INSERT INTO account (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM account
ORDER BY id ;


-- name: UpdateAccount :exec
UPDATE account 
SET balance = $2, owner = $3, currency= $4, created_at= $5
WHERE id = $1;

/*
type Account struct {
	ID        int64
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}
*/
-- name: DeleteAccount :exec
DELETE FROM account 
WHERE id = $1;

