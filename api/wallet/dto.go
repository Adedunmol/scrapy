package wallet

import (
	"github.com/google/uuid"
	"time"
)

type CreateWalletBody struct {
	Balance   float32   `json:"balance"`
	CompanyID uuid.UUID `json:"company_id"`
}

type Wallet struct {
	ID        uuid.UUID  `json:"id"`
	Balance   float32    `json:"balance"`
	CompanyID uuid.UUID  `json:"company_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
