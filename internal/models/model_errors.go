package models

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidData = errors.New("invalid data was sent")
var ErrUserNotFound = errors.New("user not found")
