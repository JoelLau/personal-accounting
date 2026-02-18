-- name: CreatePosting :one
INSERT INTO postings (
  description, 
  system_notes, 
  transacted_at
) VALUES (
  $1, $2, $3
)
RETURNING id, description, system_notes, transacted_at;


-- name: GetPosting :one
SELECT id, description, system_notes, transacted_at 
FROM postings
WHERE id = $1 LIMIT 1;


-- name: ListPostings :many
SELECT id, description, system_notes, transacted_at 
FROM postings
ORDER BY transacted_at DESC, id DESC;


-- name: UpdatePosting :one
UPDATE postings
SET
  description = $2,
  system_notes = $3,
  transacted_at = $4
WHERE id = $1
RETURNING id, description, system_notes, transacted_at;


-- name: DeletePosting :exec
DELETE FROM postings
WHERE id = $1;
