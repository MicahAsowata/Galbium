-- name: GetTodo :one
SELECT * FROM todo
WHERE id = ? LIMIT 1;

-- name: CreateTodo :execresult
INSERT INTO todo (
  name, details, completed
) VALUES (
  ?, ?, ?
);

-- name: ListTodo :many
SELECT * FROM todo
WHERE name != ''
ORDER BY created;

-- name: UpdateTodo :exec
UPDATE todo SET name = ?, details = ?, completed = ? 
WHERE id = ?;