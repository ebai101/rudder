-- name: GetAccountRows :many
select *
from accounts_view av;

-- name: GetAccountBalances :many
select av.id,
    av.account_name,
    av.balance
from accounts_view av;

-- name: GetAccount :one
select *
from accounts_view av
where av.id = $1;

-- name: GetAccountTransactions :many
select *
from transactions_view
where account_id = (
        select account_id
        from accounts a
        where a.id = $1
    )
limit $2 offset $3;

-- name: GetTransactionRows :many
select *
from transactions_view tv
where tv.description ilike $1
limit $2 offset $3;

-- name: GetTransaction :one
select *
from transactions_view tv
where tv.id = $1;