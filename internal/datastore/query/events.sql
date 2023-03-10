-- name: CreateEvent :one
INSERT INTO events (
  event_name, 
  tickets_remaining,
  event_timestamp
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE event_id = $1 LIMIT 1;

-- name: GetEvents :many
SELECT * FROM events
WHERE event_id = ANY(@ids::bigint[]);

-- name: UpdateEvent :one
UPDATE events
SET tickets_remaining = tickets_remaining - sqlc.arg(n_tickets)::int
WHERE event_id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE event_id = $1;