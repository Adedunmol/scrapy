package helpers

import (
	"github.com/google/uuid"
	"time"
)

const TokenExpiration = 30 * time.Minute

func GenerateToken(id uuid.UUID, username string) (string, error) {
	return "", nil
}
