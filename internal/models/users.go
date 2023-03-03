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
	if err != nil {
		mySQLErr := err.(*mysql.MySQLError)
		if mySQLErr.Number == 1062 {
			return errors.New("it exists already")
		}
		return err
	}
	return nil
}

const authUser = `SELECT id, password_hash FROM users WHERE email = ?`

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
			return 0, errors.New("the information you put in weren't correct")
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(arg.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, errors.New("the information you put in weren't correct")
		}
		return 0, err
	}

	return id, nil
}

func (u *Users) Exists(ctx context.Context, id int) (bool, error) {
	return false, nil
}
