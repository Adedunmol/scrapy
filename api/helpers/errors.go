package helpers

import "errors"

var ErrConflict = errors.New("user with details already exists")
