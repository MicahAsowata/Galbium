package models

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID        int64
	Name      string
	Details   string
	Completed bool
	Created   time.Time
}

type Queries struct {
	DB *sql.DB
}

const createTodo = `-- name: CreateTodo :execresult
INSERT INTO todo (
  name, details, completed
) VALUES (
  ?, ?, ?
)
`

type CreateTodoParams struct {
	Name      string
	Details   string
	Completed bool
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (sql.Result, error) {
	return q.DB.ExecContext(ctx, createTodo, arg.Name, arg.Details, arg.Completed)
}

const deleteTodo = `-- name: DeleteTodo :exec
DELETE FROM todo WHERE id = ?
`

func (q *Queries) DeleteTodo(ctx context.Context, id int64) error {
	_, err := q.DB.ExecContext(ctx, deleteTodo, id)
	return err
}

const getTodo = `-- name: GetTodo :one
SELECT id, name, details, completed, created FROM todo
WHERE id = ? LIMIT 1
`

func (q *Queries) GetTodo(ctx context.Context, id int64) (Todo, error) {
	row := q.DB.QueryRowContext(ctx, getTodo, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Details,
		&i.Completed,
		&i.Created,
	)
	return i, err
}

const listTodo = `-- name: ListTodo :many
SELECT id, name, details, completed, created FROM todo
WHERE name != ''
ORDER BY created
`

func (q *Queries) ListTodo(ctx context.Context) ([]Todo, error) {
	rows, err := q.DB.QueryContext(ctx, listTodo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Details,
			&i.Completed,
			&i.Created,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTodo = `-- name: UpdateTodo :exec
UPDATE todo SET name = ?, details = ?, completed = ? 
WHERE id = ?
`

type UpdateTodoParams struct {
	Name      string
	Details   string
	Completed bool
	ID        int64
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.DB.ExecContext(ctx, updateTodo,
		arg.Name,
		arg.Details,
		arg.Completed,
		arg.ID,
	)
	return err
}
