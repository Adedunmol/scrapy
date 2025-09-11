package wallet

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type CreateWalletBody struct {
	CompanyID uuid.UUID `json:"company_id"`
}

type Wallet struct {
	ID        uuid.UUID       `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	CompanyID uuid.UUID       `json:"company_id"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}
