-- name: InsertOrganizations :batchexec
insert into organizations (inst_name, sfin_url, domain_name, url)
values ($1, $2, $3, $4) on conflict do nothing;

-- name: InsertAccounts :batchexec
insert into accounts (account_id, account_name, inst_name, currency)
values ($1, $2, $3, $4) on conflict do nothing;

-- name: InsertBalances :batchexec
insert into balances (balance_id, balance_date, balance, account_id)
values ($1, $2, $3, $4) on conflict do nothing;

-- name: InsertTransactions :batchexec
insert into transactions (
        transaction_id,
        posted_date,
        description,
        amount,
        account_id,
        inst_name,
        full_description
    )
values ($1, $2, $3, $4, $5, $6, $7) on conflict do nothing;