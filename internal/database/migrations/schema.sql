create type public."account_class_t" as ENUM ('Asset', 'Liability');

create type public."account_type_t" as ENUM ('CHECKING', 'SAVINGS', 'CREDIT', 'VENMO');

create type public."autocat_criteria_operator_t" as ENUM (
	'equals',
	'contains',
	'starts_with',
	'ends_with',
	'regex',
	'polarity',
	'min',
	'max'
);

-- public.autocat_rules definition
create table autocat_rules (
	id bigserial not null,
	priority int4 default 0 not null,
	created_at timestamptz default current_timestamp null,
	updated_at timestamptz default current_timestamp null,
	constraint autocat_rules_pkey primary key (id)
);

-- public.categories definition
create table categories (
	id bigserial not null,
	category text not null,
	category_group text not null,
	category_type text not null,
	business_expense bool default false not null,
	hide_from_reports bool default false not null,
	constraint categories_pk primary key (id),
	constraint categories_unique unique (category)
);

-- public.organizations definition
create table organizations (
	id bigserial not null,
	inst_name text not null,
	sfin_url text not null,
	domain_name text null,
	url text null,
	constraint organizations_pk primary key (id),
	constraint organizations_unique unique (inst_name)
);

-- public.transactions_cols definition
create table transactions_cols (
	column_name text not null,
	data_type text not null,
	constraint transactions_cols_pk primary key (column_name)
);

-- public.accounts definition
create table accounts (
	id bigserial not null,
	account_id text not null,
	account_name text not null,
	inst_name text not null,
	account_type public."account_type_t" null,
	account_class public."account_class_t" null,
	currency text not null,
	active bool default true not null,
	constraint accounts_pk primary key (id),
	constraint accounts_unique unique (account_id),
	constraint accounts_organizations_fk foreign key (inst_name) references organizations (inst_name) on delete restrict on update cascade
);

-- public.autocat_criteria definition
create table autocat_criteria (
	id bigserial not null,
	rule_id int4 not null,
	column_name varchar(255) not null,
	"operator" public."autocat_criteria_operator_t" not null,
	filter_value_text text null,
	filter_value_numeric numeric(10, 2) null,
	filter_value_timestamptz timestamptz null,
	criteria_order int4 null,
	constraint rule_criteria_pkey primary key (id),
	constraint uq_rule_criteria_order unique (rule_id, criteria_order),
	constraint autocat_criteria_autocat_rules_fk foreign key (rule_id) references autocat_rules (id) on delete cascade,
	constraint autocat_criteria_transactions_cols_fk foreign key (column_name) references transactions_cols (column_name)
);

-- Table Triggers
create trigger validate_filter_value_trigger before
insert
	or
update on public.autocat_criteria for each row execute function validate_filter_value ();

-- public.autocat_overrides definition
create table autocat_overrides (
	id bigserial not null,
	rule_id int8 not null,
	column_name varchar(255) not null,
	override_value text not null,
	override_order int4 not null,
	constraint rule_overrides_pkey primary key (id),
	constraint uq_rule_overrides_order unique (rule_id, override_order),
	constraint autocat_overrides_autocat_rules_fk foreign key (rule_id) references autocat_rules (id) on delete cascade
);

-- public.balances definition
create table balances (
	id bigserial not null,
	balance_id text not null,
	balance_date timestamptz not null,
	balance numeric(10, 2) not null,
	account_id text not null,
	added_date timestamptz default now () not null,
	constraint balances_pk primary key (id),
	constraint balances_unique unique (balance_id),
	constraint balances_accounts_fk foreign key (account_id) references accounts (account_id) on delete restrict on update cascade
);

-- public.transactions definition
create table transactions (
	id bigserial not null,
	transaction_id text not null,
	posted_date timestamptz not null,
	description text null,
	category text null,
	amount numeric(10, 2) not null,
	account_id text not null,
	inst_name text not null,
	full_description text not null,
	added_date timestamptz default now () not null,
	categorized_date timestamptz null,
	note text null,
	check_num text null,
	constraint transactions_pk primary key (id),
	constraint transactions_unique unique (transaction_id),
	constraint transactions_accounts_fk foreign key (account_id) references accounts (account_id) on delete restrict on update cascade,
	constraint transactions_categories_fk foreign key (category) references categories (category) on delete restrict on update cascade,
	constraint transactions_organizations_fk foreign key (inst_name) references organizations (inst_name) on delete restrict on update cascade
);

-- public.transactions_view source
create or replace view transactions_view as
select t.id,
	t.transaction_id,
	t.posted_date,
	t.description,
	t.category,
	t.amount,
	t.account_id,
	a.account_name,
	t.inst_name,
	t.full_description,
	t.added_date,
	t.categorized_date
from transactions t
	join accounts a on t.account_id = a.account_id
order by t.posted_date desc;

-- public.accounts_view source
create or replace view accounts_view as with ranked_balances as (
		select balances.account_id,
			balances.balance,
			balances.balance_date,
			balances.added_date,
			row_number() OVER (
				partition BY balances.account_id
				order by balances.balance_date desc
			) as rank
		from balances
	)
select a.id,
	a.account_id,
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