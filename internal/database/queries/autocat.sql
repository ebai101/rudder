-- name: MatchTransactions :many
with rule_criteria as (
    select r.id as rule_id,
        jsonb_agg(
            jsonb_build_object(
                'column_name',
                c.column_name,
                'operator',
                c.operator::text,
                'filter_value_text',
                c.filter_value_text,
                'filter_value_numeric',
                c.filter_value_numeric,
                'filter_value_timestamptz',
                c.filter_value_timestamptz
            )
        ) as criteria
    from autocat_rules r
        join autocat_criteria c on r.id = c.rule_id
    group by r.id
),
rule_overrides as (
    select rule_id,
        jsonb_agg(
            jsonb_build_object(
                'column_name',
                column_name,
                'override_value',
                override_value,
                'override_order',
                override_order
            )
            order by override_order
        ) as overrides
    from autocat_overrides
    group by rule_id
),
matching_transactions as (
    select t.*,
        rc.rule_id,
        rc.criteria
    from transactions t
        cross join rule_criteria rc
    where (
            select bool_and(
                    COALESCE(
                        check_transaction_criteria(t, c),
                        false
                    )
                )
            from jsonb_array_elements(to_jsonb(criteria)) as c
        )
        and t.category is null
),
processed_overrides as (
    select rule_id,
        MAX(
            case
                when column_name = 'category' then override_value
            end
        ) as category_override
    from (
            select rule_id,
                jsonb_array_elements(overrides)->>'column_name' as column_name,
                jsonb_array_elements(overrides)->>'override_value' as override_value
            from rule_overrides
        ) as override_exploded
    group by rule_id
)
select mt.transaction_id,
    mt.rule_id,
    COALESCE(
        po.category_override,
        mt.category
    ) as new_category
from matching_transactions mt
    left join rule_overrides ro on mt.rule_id = ro.rule_id
    left join processed_overrides po on mt.rule_id = po.rule_id
order by mt.rule_id,
    mt.id;

-- name: UpdateTransactionCategories :batchexec
update transactions
set category = $2,
    categorized_date = $3
where transaction_id = $1;

-- name: GetAutocatRules :many
select r.id,
    (
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