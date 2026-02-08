-- name: CreateLedgerAccount :one
INSERT INTO ledger_accounts (
  name,
  description,
  parent_id
) VALUES (
  $1, $2, $3
)
RETURNING id, name, description, parent_id;


-- name: GetLedgerAccount :one
SELECT id, name, description, parent_id 
FROM ledger_accounts
WHERE id = $1 LIMIT 1;


-- name: ListLedgerAccounts :many
WITH RECURSIVE account_tree AS (
    -- Base Case: Select root accounts
    SELECT 
        id,
        name,
        description,
        parent_id,
        name::TEXT AS qualified_name,
        1 AS level
    FROM ledger_accounts
    WHERE parent_id IS NULL

    UNION ALL

    -- Recursive Step: Join children to parents
    SELECT 
        child.id,
        child.name,
        child.description,
        child.parent_id,
        (parent.qualified_name || ':' || child.name)::TEXT AS qualified_name,
        parent.level + 1 AS level
    FROM ledger_accounts child
    JOIN account_tree parent ON child.parent_id = parent.id
)
SELECT id, name, qualified_name, description, parent_id, level
FROM account_tree 
ORDER BY qualified_name;


-- name: UpdateLedgerAccount :one
UPDATE ledger_accounts
SET
  name = $2,
  description = $3,
  parent_id = $4
WHERE id = $1
RETURNING id, name, description, parent_id;


-- name: DeleteLedgerAccount :exec
DELETE FROM ledger_accounts
WHERE id = $1;
