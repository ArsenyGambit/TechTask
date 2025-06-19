package errors

import "errors"

var (
	ErrNewsNotFound      = errors.New("news not found")
	ErrDuplicateSlug     = errors.New("news with this slug already exists")
	ErrInvalidSlug       = errors.New("invalid slug format")
	ErrInvalidTitle      = errors.New("invalid title")
	ErrInvalidContent    = errors.New("invalid content")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)
