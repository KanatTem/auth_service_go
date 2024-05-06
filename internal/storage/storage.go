package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
	ErrUserNotAdmin = errors.New("user is not Admin")
	ErrRoleNotFound = errors.New("role not found")
	ErrUserHaveRole = errors.New("user alredy have role")
)
