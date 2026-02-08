-- name: CreateEntry :one
INSERT INTO entries (
  description
  , system_notes
  , postings_id
  , ledger_accounts_id
  , debit_microsgd
  , credit_microsgd
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, description, system_notes, postings_id, ledger_accounts_id, debit_microsgd, credit_microsgd;


-- name: GetEntry :one
SELECT id, description, system_notes, postings_id, ledger_accounts_id, debit_microsgd, credit_microsgd
FROM entries
WHERE id = $1 LIMIT 1;


-- name: ListEntries :many
SELECT id, description, system_notes, postings_id, ledger_accounts_id, debit_microsgd, credit_microsgd
FROM entries;


-- name: UpdateEntry :one
UPDATE entries
SET
  description = $2
  , system_notes = $3
  , ledger_accounts_id = $4
  , debit_microsgd = $5
  , credit_microsgd = $6
WHERE id = $1
RETURNING id, description, system_notes, postings_id, ledger_accounts_id, debit_microsgd, credit_microsgd;


-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
