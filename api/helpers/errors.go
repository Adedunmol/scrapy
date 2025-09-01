package helpers

import "errors"

var ErrConflict = errors.New("user with details already exists")
var ErrNotFound = errors.New("resource not found")
var ErrInternalServer = errors.New("internal server error")
