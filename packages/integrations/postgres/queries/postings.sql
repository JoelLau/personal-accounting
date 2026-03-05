-- name: CreatePosting :one
INSERT INTO postings(description, system_notes, transacted_at, source_hash)
  VALUES ($1, $2, $3, $4)
RETURNING
  id, description, system_notes, transacted_at, source_hash;

-- name: GetPosting :one
SELECT
  id,
  description,
  system_notes,
  transacted_at,
  source_hash
FROM
  postings
WHERE
  id = $1
LIMIT 1;

-- name: ListPostings :many
SELECT
  id,
  description,
  system_notes,
  transacted_at,
  source_hash
FROM
  postings
ORDER BY
  transacted_at DESC,
  id DESC;

-- name: UpdatePosting :one
UPDATE
  postings
SET
  description = $2,
  system_notes = $3,
  transacted_at = $4,
  source_hash = $5
WHERE
  id = $1
RETURNING
  id,
  description,
  system_notes,
  transacted_at,
  source_hash;

-- name: DeletePosting :exec
DELETE FROM postings
WHERE id = $1;

