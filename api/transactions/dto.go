package transactions

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type CreateTransactionBody struct {
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	Reference     string          `json:"reference"`
	Status        string          `json:"status"`
	Type          string          `json:"type"`
	WalletID      uuid.UUID       `json:"wallet_id"`
}

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	Reference     string          `json:"reference"`
	Status        string          `json:"status"`
	Type          string          `json:"type"`
	WalletID      uuid.UUID       `json:"wallet_id"`
	CreatedAt     *time.Time      `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}
