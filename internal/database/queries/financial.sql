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
with current_assets as (
    select COALESCE(SUM(balance)::numeric, 0::numeric) as current_assets
    from accounts_view av
    where av.account_class = 'Asset'
),
current_liabilities as (
    select ABS(COALESCE(SUM(balance)::numeric, 0::numeric)) as current_liabilities
    from accounts_view av
    where av.account_class = 'Liability'
),
total_income as (
    select coalesce(sum(amount)::numeric, 0::numeric) as total_income
    from transactions_view tv
    where amount > 0
	and tv.posted_date <= $2
	and tv.posted_date >= $1
),
total_expense as (
    select ABS(COALESCE(SUM(amount)::numeric, 0::numeric)) as total_expense
    from transactions_view tv
    where amount < 0
	and tv.posted_date <= $2
	and tv.posted_date >= $1
),
needs_cat as (
	select
		count(*) as needs_cat_count,
		coalesce(sum(amount)::numeric, 0::numeric) as needs_cat_amt
	from transactions_view
	where category is null
)
select total_income::numeric,
	total_expense::numeric,
	(total_income - total_expense)::numeric as cash_flow,
    current_assets::numeric,
    current_liabilities::numeric,
    (current_assets - current_liabilities)::numeric as net_worth,
    needs_cat_count,
    needs_cat_amt::numeric
from total_income
	cross join total_expense
    cross join current_assets
    cross join current_liabilities
	cross join needs_cat;

