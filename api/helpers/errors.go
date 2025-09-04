package helpers

import "errors"

var ErrConflict = errors.New("record already exists")
var ErrNotFound = errors.New("resource not found")
var ErrInternalServer = errors.New("internal server error")
