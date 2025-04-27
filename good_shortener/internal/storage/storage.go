package storage

import "errors"

var (
	ErrURLExists   = errors.New("url exists")
	ErrURLNotFound = errors.New("url not found")
	UserNotFound   = errors.New("user not found")
	UserExists     = errors.New("user exist")
)
