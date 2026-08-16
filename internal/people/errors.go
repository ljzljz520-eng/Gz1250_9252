package people

import "errors"

var (
	ErrContactConflict = errors.New("contact conflict")
	ErrNotFound        = errors.New("person not found")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidStatus   = errors.New("invalid status")
)
