package models

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID        int
	Name      string
	Details   string
	Completed bool
	Created   time.Time
}

type Queries struct {
	DB *sql.DB
}

const createTodo = `INSERT INTO todo (name, details, completed, user_id) 
VALUES (?, ?, ?, ?)`

type CreateTodoParams struct {
	Name      string
	Details   string
	Completed bool
	UserID    int
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (sql.Result, error) {
	return q.DB.ExecContext(ctx, createTodo, arg.Name, arg.Details, arg.Completed, arg.UserID)
}

const deleteTodo = `DELETE FROM todo WHERE id = ? AND user_id = ?`

func (q *Queries) DeleteTodo(ctx context.Context, id, user_id int) error {
	_, err := q.DB.ExecContext(ctx, deleteTodo, id, user_id)
	return err
}

const getTodo = `SELECT id, name, details, completed, created FROM todo
WHERE id = ? AND user_id = ? LIMIT 1`

func (q *Queries) GetTodo(ctx context.Context, id, user_id int) (Todo, error) {
	row := q.DB.QueryRowContext(ctx, getTodo, id, user_id)
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

const listTodo = `SELECT id, name, details, completed, created FROM todo
WHERE name != '' AND user_id = ?
ORDER BY created`

func (q *Queries) ListTodo(ctx context.Context, userID int) ([]Todo, error) {
	rows, err := q.DB.QueryContext(ctx, listTodo, userID)
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

const updateTodo = `UPDATE todo SET name = ?, details = ?, completed = ? WHERE id = ? AND user_id = ?`

type UpdateTodoParams struct {
	Name      string
	Details   string
	Completed bool
	ID        int
	UserID    int
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.DB.ExecContext(ctx, updateTodo,
		arg.Name,
		arg.Details,
		arg.Completed,
		arg.ID,
		arg.UserID,
	)
	return err
}
