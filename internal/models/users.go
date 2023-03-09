package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

const insertUser = `INSERT INTO users ( name, email, username, password_hash)
	VALUES (?, ?, ?, ?)`

type InsertUsersParams struct {
	Name     string
	Email    string
	Username string
	Password string
}

func (u *Users) Insert(ctx context.Context, arg InsertUsersParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), 12)
	if err != nil {
		return err
	}
	_, err = u.DB.ExecContext(ctx, insertUser, arg.Name, arg.Email, arg.Username, hashedPassword)
	var mySQLError mysql.MySQLError
	if err != nil {
		if errors.Is(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				return ErrUserAlreadyExists
			}
		}
		return err
	}
	return nil
}

const authUser = `SELECT user_id, password_hash FROM users WHERE email = ?`

type AuthUserParams struct {
	Email    string
	Password string
}

func (u *Users) Authenticate(ctx context.Context, arg AuthUserParams) (int, error) {
	var id int
	var passwordHash []byte

	err := u.DB.QueryRowContext(ctx, authUser, arg.Email).Scan(&id, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidData
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(arg.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidData
		}
		return 0, err
	}

	return id, nil
}

const getUserName = `SELECT username FROM users WHERE user_id = ?`

func (u *Users) Get(ctx context.Context, id int) (string, error) {
	var username string
	err := u.DB.QueryRowContext(ctx, getUserName, id).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidData
		}
		return "", err
	}
	return username, nil
}

const resetPassword = `UPDATE users SET password_hash = ? WHERE email = ?`
const findAccount = `SELECT EXISTS(SELECT email FROM users WHERE email = ?)`

type ResetPasswordParams struct {
	Email    string
	Password string
}

func (u *Users) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	var exists bool

	err := u.DB.QueryRowContext(ctx, findAccount, arg.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), 12)
	if err != nil {
		return err
	}
	err = u.DB.QueryRowContext(ctx, resetPassword, arg.Email, hashedPassword).Err()
	if err != nil {
		return err
	}
	return nil
}
