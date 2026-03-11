-- name: GetMonthlyAccountRollup :many
WITH RECURSIVE
-- 1. Create a list of every month in the requested range
month_range AS (
  SELECT
    generate_series(date_trunc('month', @start_date::timestamptz), date_trunc('month', @end_date::timestamptz), '1 month'::interval)::timestamptz AS month_date
),
-- 2. Build the qualified names for all accounts
account_paths AS (
  SELECT
    id,
    parent_id,
    name AS qualified_name
  FROM
    ledger_accounts
  WHERE
    parent_id IS NULL
  UNION ALL
  SELECT
    child.id,
    child.parent_id,
    (parent.qualified_name || ':' || child.name)
  FROM
    ledger_accounts child
    JOIN account_paths parent ON child.parent_id = parent.id
),
-- 3. Calculate raw activity filtered by date
base_activity AS (
  SELECT
    e.ledger_accounts_id AS id,
    date_trunc('month', p.transacted_at)::timestamptz AS month_date,
    SUM(e.debit_microsgd)::bigint AS ind_debit,
    SUM(e.credit_microsgd)::bigint AS ind_credit
  FROM
    entries e
    JOIN postings p ON e.postings_id = p.id
  WHERE
    p.transacted_at >= @start_date
    AND p.transacted_at <= @end_date
  GROUP BY
    e.ledger_accounts_id,
    date_trunc('month', p.transacted_at)
),
-- 4. Propagate totals up the hierarchy
rolled_up_totals AS (
  SELECT
    id,
    month_date,
    ind_debit AS roll_debit,
    ind_credit AS roll_credit
  FROM
    base_activity
  UNION ALL
  SELECT
    la.parent_id,
    rt.month_date,
    rt.roll_debit,
    rt.roll_credit
  FROM
    rolled_up_totals rt
    JOIN ledger_accounts la ON rt.id = la.id
  WHERE
    la.parent_id IS NOT NULL
)
-- 5. Final Join: Every account x Every month
SELECT
  ap.id AS ledger_account_id,
  ap.parent_id,
  ap.qualified_name,
  mr.month_date,
  COALESCE(ba.ind_debit, 0)::bigint AS individual_debit_microsgd,
  COALESCE(ba.ind_credit, 0)::bigint AS individual_credit_microsgd,
  COALESCE(SUM(rt.roll_debit), 0)::bigint AS rolled_up_debit_microsgd,
  COALESCE(SUM(rt.roll_credit), 0)::bigint AS rolled_up_credit_microsgd,
  (COALESCE(SUM(rt.roll_debit), 0) - COALESCE(SUM(rt.roll_credit), 0))::bigint AS rolled_up_net_microsgd
FROM
  month_range mr
CROSS JOIN
  account_paths ap
LEFT JOIN
  base_activity ba ON ba.id = ap.id AND ba.month_date = mr.month_date
LEFT JOIN
  rolled_up_totals rt ON rt.id = ap.id AND rt.month_date = mr.month_date
GROUP BY
  ap.id,
  ap.parent_id,
  ap.qualified_name,
  mr.month_date,
  ba.ind_debit,
  ba.ind_credit
ORDER BY
  mr.month_date DESC,
  ap.qualified_name ASC;
