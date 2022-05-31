package todo

import "errors"

var (
	ErrProjectNotFound = errors.New("todo: project not found")
	ErrVersionMismatch = errors.New("todo: project version mismatch")
	ErrIDsMismatch     = errors.New("todo: project ids mismatch")
	ErrAlreadyExists   = errors.New("todo: project already exists")
)
