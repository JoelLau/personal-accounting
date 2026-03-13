-- name: GetUncategorizedExpenses :many
WITH accounts AS (WITH RECURSIVE account_tree AS (
    -- Base Case: Select root accounts
    SELECT id,
           name,
           description,
           parent_id,
           name::text AS qualified_name,
           1          AS level
    FROM ledger_accounts
    WHERE parent_id IS NULL
    UNION ALL
    -- Recursive Step: Join children to parents
    SELECT child.id,
           child.name,
           child.description,
           child.parent_id,
           (parent.qualified_name || ':' || child.name)::text AS qualified_name,
           parent.level + 1                                   AS level
    FROM ledger_accounts child
             JOIN account_tree parent ON child.parent_id = parent.id)
                  SELECT id,
                         name,
                         qualified_name,
                         description,
                         parent_id,
                         level
                  FROM account_tree
                  ORDER BY qualified_name)
SELECT
--     p.id
    e.id
     , TO_CHAR(p.transacted_at, 'yyyy-mm-dd') as transacted_at
     , p.description
     , (e.debit_microsgd / 1_000_000.0)       as debit
     , (e.credit_microsgd / 1_000_000.0)      as credit
     , e.ledger_accounts_id
     , acc.qualified_name                     as category
FROM entries e
         LEFT JOIN public.postings p on e.postings_id = p.id
         LEFT JOIN accounts acc on e.ledger_accounts_id = acc.id
WHERE e.ledger_accounts_id = 4000
  AND e.debit_microsgd > 0
  AND e.credit_microsgd <= 0
  AND p.transacted_at BETWEEN '2026-02-01' AND '2026-02-28'
ORDER BY debit_microsgd DESC
       , description DESC;

-- name: GetUncategorizedRefunds :many
WITH accounts AS (WITH RECURSIVE account_tree AS (
    -- Base Case: Select root accounts
    SELECT id,
           name,
           description,
           parent_id,
           name::text AS qualified_name,
           1          AS level
    FROM ledger_accounts
    WHERE parent_id IS NULL
    UNION ALL
    -- Recursive Step: Join children to parents
    SELECT child.id,
           child.name,
           child.description,
           child.parent_id,
           (parent.qualified_name || ':' || child.name)::text AS qualified_name,
           parent.level + 1                                   AS level
    FROM ledger_accounts child
             JOIN account_tree parent ON child.parent_id = parent.id)
                  SELECT id,
                         name,
                         qualified_name,
                         description,
                         parent_id,
                         level
                  FROM account_tree
                  ORDER BY qualified_name)
SELECT
--     p.id
    e.id
     , TO_CHAR(p.transacted_at, 'yyyy-mm-dd') as transacted_at
     , p.description
     , (e.debit_microsgd / 1_000_000.0)       as debit
     , (e.credit_microsgd / 1_000_000.0)      as credit
     , e.ledger_accounts_id
     , acc.qualified_name                     as category
FROM entries e
         LEFT JOIN public.postings p on e.postings_id = p.id
         LEFT JOIN accounts acc on e.ledger_accounts_id = acc.id
WHERE e.ledger_accounts_id = 4000
  AND e.debit_microsgd <= 0
  AND e.credit_microsgd > 0
ORDER BY transacted_at DESC
       , description DESC;

-- name: ListDebits :many
WITH accounts AS (WITH RECURSIVE account_tree AS (
    -- Base Case: Select root accounts
    SELECT id,
           name,
           description,
           parent_id,
           name::text AS qualified_name,
           1          AS level
    FROM ledger_accounts
    WHERE parent_id IS NULL
    UNION ALL
    -- Recursive Step: Join children to parents
    SELECT child.id,
           child.name,
           child.description,
           child.parent_id,
           (parent.qualified_name || ':' || child.name)::text AS qualified_name,
           parent.level + 1                                   AS level
    FROM ledger_accounts child
             JOIN account_tree parent ON child.parent_id = parent.id)
                  SELECT id,
                         name,
                         qualified_name,
                         description,
                         parent_id,
                         level
                  FROM account_tree
                  ORDER BY qualified_name)
SELECT
--     p.id
--     e.id
     TO_CHAR(p.transacted_at, 'yyyy-mm-dd') as transacted_at
     , p.description
     , (e.debit_microsgd / 1_000_000.0)       as debit
     , (e.credit_microsgd / 1_000_000.0)      as credit
--      , e.ledger_accounts_id
     , acc.qualified_name                     as category
FROM entries e
         LEFT JOIN public.postings p on e.postings_id = p.id
         LEFT JOIN accounts acc on e.ledger_accounts_id = acc.id
WHERE e.debit_microsgd > 0
  AND e.credit_microsgd <= 0
ORDER BY e.debit_microsgd DESC, transacted_at DESC, description DESC
