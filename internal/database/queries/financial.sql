-- name: GetAccountRows :many
with ranked_balances as (
    select account_id,
        balance,
        balance_date,
        added_date,
        row_number() over (
            partition by account_id
            order by balance_date desc
        ) as rank
    from balances
)
select a.account_id,
    a.account_name,
    a.inst_name,
    a.account_type,
    a.account_class,
    a.currency,
    a.active,
    rb.balance,
    rb.balance_date,
    rb.added_date
from accounts a
    join ranked_balances rb on a.account_id = rb.account_id
    and rb.rank = 1;

-- name: GetTransactionRows :many
select t.transaction_id,
    t.posted_date,
    t.description,
    t.category,
    t.amount,
    a.account_name,
    t.inst_name,
    t.full_description,
    t.added_date,
    t.categorized_date,
    t.note,
    t.check_num
from transactions t
    join accounts a on t.account_id = a.account_id
order by posted_date desc
limit $1;