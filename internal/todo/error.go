package todo

import (
	"errors"
	"strings"
)

var (
	ErrProjectNotFound = errors.New("todo: project not found")
	ErrVersionMismatch = errors.New("todo: project version mismatch")
	ErrIDsMismatch     = errors.New("todo: project ids mismatch")
	ErrAlreadyExists   = errors.New("todo: project already exists")
)

func IsStorageError(err error) bool {
	if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrVersionMismatch) ||
		errors.Is(err, ErrIDsMismatch) || errors.Is(err, ErrAlreadyExists) {
		return true
	}

	return false
}

func IsDuplicateKeyError(err error) bool {
	if strings.Contains(err.Error(), "duplicate key error") {
		return true
	}

	return false
}
