-- name: GetAutocatRules :many
select (
        select json_agg (
                json_build_object (
                    'column_name',
                    c.column_name,
                    'operator',
                    c."operator",
                    'filter_value_text',
                    c.filter_value_text,
                    'filter_value_numeric',
                    c.filter_value_numeric,
                    'filter_value_timestamptz',
                    c.filter_value_timestamptz,
                    'criteria_order',
                    c.criteria_order
                )
            )
        from autocat_criteria c
        where c.rule_id = r.id
    ) as criteria,
    (
        select json_agg (
                json_build_object (
                    'column_name',
                    o.column_name,
                    'override_value',
                    o.override_value,
                    'override_order',
                    o.override_order
                )
            )
        from autocat_overrides o
        where o.rule_id = r.id
    ) as overrides
from autocat_rules r;
-- name: InsertOrganizations :one
with rows as (
    insert into organizations (inst_name, sfin_url, domain_name, url)
    values ($1, $2, $3, $4) on conflict do nothing
    returning 1
)
select count(*)
from rows;
-- name: InsertAccounts :one
with rows as (
    insert into accounts (account_id, account_name, inst_name, currency)
    values ($1, $2, $3, $4) on conflict do nothing
    returning 1
)
select count(*)
from rows;
-- name: InsertBalances :one
with rows as (
    insert into balances (balance_id, balance_date, balance, account_id)
    values ($1, $2, $3, $4) on conflict do nothing
    returning 1
)
select count(*)
from rows;
-- name: InsertTransactions :one
with rows as (
    insert into transactions (
            transaction_id,
            posted_date,
            description,
            amount,
            account_id,
            inst_name,
            full_description
        )
    values ($1, $2, $3, $4, $5, $6, $7) on conflict do nothing
    returning 1
)
select count(*)
from rows;
-- name: UpdateTransactionCategories :one
with rows as (
    update transactions
    set category = $2,
        categorized_date = $3
    where transaction_id = $1
        and category is null
    returning 1
)
select count(*)
from rows;
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