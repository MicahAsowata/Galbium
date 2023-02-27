-- name: GetTodo :one
SELECT * FROM todo
WHERE id = ? AND completed != TRUE LIMIT 1;

-- name: CreateTodo :execresult
INSERT INTO todo (
  name, details, completed
) VALUES (
  ?, ?, ?
);