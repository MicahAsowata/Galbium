package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Name         string
	Email        string
	Username     string
	PasswordHash []byte
	Created      time.Time
}

type Users struct {
	DB *sql.DB
}

const insertUser string = ""

type InsertUsersParams struct {
	Name     string
	Email    string
	Username string
	Password string
}

func (u *Users) Insert(ctx context.Context, arg InsertUsersParams) error {
	return nil
}

type AuthUserParams struct {
	Email    string
	Password string
}

func (u *Users) Authenticate(ctx context.Context, arg AuthUserParams) (int, error) {
	return 0, nil
}

func (u *Users) Exists(ctx context.Context, id int) (bool, error) {
	return false, nil
}
