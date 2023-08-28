-- name: CreateEntry :one
INSERT INTO entry (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;