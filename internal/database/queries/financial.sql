-- name: GetOrganizationRows :many
select inst_name,
    sfin_url,
    domain_name,
    url
from organizations;

-- name: GetAccountRows :many
select account_id,
    account_name,
    inst_name,
    account_type,
    account_class,
    currency,
    active
from accounts;

-- name: GetBalanceRows :many
select balance_id,
    balance_date,
    balance,
    account_id,
    added_date
from balances;

-- name: GetTransactionRows :many
select transaction_id,
    posted_date,
    description,
    category,
    amount,
    account_id,
    inst_name,
    full_description,
    added_date,
    categorized_date,
    note,
    check_num
from transactions;