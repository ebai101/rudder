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

-- name: GetInsights :one
with spent_week as (
    select ABS(COALESCE(SUM(amount)::numeric, 0::numeric)) as spent_week
    from transactions_view
    where amount < 0
        and posted_date >= current_date AT TIME ZONE 'UTC' - interval '7 days'
),
total_assets as (
    select COALESCE(SUM(balance)::numeric, 0::numeric) as total_assets
    from accounts_view
    where account_class = 'Asset'
),
total_liabilities as (
    select ABS(COALESCE(SUM(balance)::numeric, 0::numeric)) as total_liabilities
    from accounts_view
    where account_class = 'Liability'
)
select spent_week::numeric,
    total_assets::numeric,
    total_liabilities::numeric,
    (total_assets - total_liabilities)::numeric as net_worth
from spent_week
    cross join total_assets
    cross join total_liabilities;