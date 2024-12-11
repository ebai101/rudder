CREATE TYPE public."account_class_t" AS ENUM ('Asset', 'Liability');
CREATE TYPE public."account_type_t" AS ENUM ('CHECKING', 'SAVINGS', 'CREDIT', 'VENMO');
CREATE TYPE public."autocat_criteria_operator_t" AS ENUM (
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
CREATE TABLE autocat_rules (
	id bigserial NOT NULL,
	priority int4 DEFAULT 0 NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT autocat_rules_pkey PRIMARY KEY (id)
);
-- public.categories definition
CREATE TABLE categories (
	id bigserial NOT NULL,
	category text NOT NULL,
	category_group text NOT NULL,
	category_type text NOT NULL,
	business_expense bool DEFAULT false NOT NULL,
	hide_from_reports bool DEFAULT false NOT NULL,
	CONSTRAINT categories_pk PRIMARY KEY (id),
	CONSTRAINT categories_unique UNIQUE (category)
);
-- public.organizations definition
CREATE TABLE organizations (
	id bigserial NOT NULL,
	inst_name text NOT NULL,
	sfin_url text NOT NULL,
	domain_name text NULL,
	url text NULL,
	CONSTRAINT organizations_pk PRIMARY KEY (id),
	CONSTRAINT organizations_unique UNIQUE (inst_name)
);
-- public.transactions_cols definition
CREATE TABLE transactions_cols (
	column_name text NOT NULL,
	data_type text NOT NULL,
	CONSTRAINT transactions_cols_pk PRIMARY KEY (column_name)
);
-- public.accounts definition
CREATE TABLE accounts (
	id bigserial NOT NULL,
	account_id text NOT NULL,
	account_name text NOT NULL,
	inst_name text NOT NULL,
	account_type public."account_type_t" NULL,
	account_class public."account_class_t" NULL,
	currency text NOT NULL,
	active bool DEFAULT true NOT NULL,
	CONSTRAINT accounts_pk PRIMARY KEY (id),
	CONSTRAINT accounts_unique UNIQUE (account_id),
	CONSTRAINT accounts_organizations_fk FOREIGN KEY (inst_name) REFERENCES organizations (inst_name) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- public.autocat_criteria definition
CREATE TABLE autocat_criteria (
	id bigserial NOT NULL,
	rule_id int4 NOT NULL,
	column_name varchar(255) NOT NULL,
	"operator" public."autocat_criteria_operator_t" NOT NULL,
	filter_value_text text NULL,
	filter_value_numeric numeric(10, 2) NULL,
	filter_value_timestamptz timestamptz NULL,
	criteria_order int4 NULL,
	CONSTRAINT rule_criteria_pkey PRIMARY KEY (id),
	CONSTRAINT uq_rule_criteria_order UNIQUE (rule_id, criteria_order),
	CONSTRAINT autocat_criteria_autocat_rules_fk FOREIGN KEY (rule_id) REFERENCES autocat_rules (id) ON DELETE CASCADE,
	CONSTRAINT autocat_criteria_transactions_cols_fk FOREIGN KEY (column_name) REFERENCES transactions_cols (column_name)
);
-- Table Triggers
create trigger validate_filter_value_trigger before
insert
	or
update on public.autocat_criteria for each row execute function validate_filter_value ();
-- public.autocat_overrides definition
CREATE TABLE autocat_overrides (
	id bigserial NOT NULL,
	rule_id int8 NOT NULL,
	column_name varchar(255) NOT NULL,
	override_value text NOT NULL,
	override_order int4 NOT NULL,
	CONSTRAINT rule_overrides_pkey PRIMARY KEY (id),
	CONSTRAINT uq_rule_overrides_order UNIQUE (rule_id, override_order),
	CONSTRAINT autocat_overrides_autocat_rules_fk FOREIGN KEY (rule_id) REFERENCES autocat_rules (id) ON DELETE CASCADE
);
-- public.balances definition
CREATE TABLE balances (
	id bigserial NOT NULL,
	balance_id text NOT NULL,
	balance_date timestamptz NOT NULL,
	balance numeric(10, 2) NOT NULL,
	account_id text NOT NULL,
	added_date timestamptz DEFAULT now () NOT NULL,
	CONSTRAINT balances_pk PRIMARY KEY (id),
	CONSTRAINT balances_unique UNIQUE (balance_id),
	CONSTRAINT balances_accounts_fk FOREIGN KEY (account_id) REFERENCES accounts (account_id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- public.transactions definition
CREATE TABLE transactions (
	id bigserial NOT NULL,
	transaction_id text NOT NULL,
	posted_date timestamptz NOT NULL,
	description text NULL,
	category text NULL,
	amount numeric(10, 2) NOT NULL,
	account_id text NOT NULL,
	inst_name text NOT NULL,
	full_description text NOT NULL,
	added_date timestamptz DEFAULT now () NOT NULL,
	categorized_date timestamptz NULL,
	note text NULL,
	check_num text NULL,
	CONSTRAINT transactions_pk PRIMARY KEY (id),
	CONSTRAINT transactions_unique UNIQUE (transaction_id),
	CONSTRAINT transactions_accounts_fk FOREIGN KEY (account_id) REFERENCES accounts (account_id) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT transactions_categories_fk FOREIGN KEY (category) REFERENCES categories (category) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT transactions_organizations_fk FOREIGN KEY (inst_name) REFERENCES organizations (inst_name) ON DELETE RESTRICT ON UPDATE CASCADE
);