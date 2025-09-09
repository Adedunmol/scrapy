package wallet

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	CreateWallet(ctx context.Context, body *CreateWalletBody) error
	GetWallet(ctx context.Context, walletID uuid.UUID) (Wallet, error)
	TopUpWallet(ctx context.Context, walletID uuid.UUID, amount float32) error
	ChargeWallet(ctx context.Context, walletID uuid.UUID, amount float32) error
}

type WalletStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewWalletStore(db *pgxpool.Pool, queryTimeout time.Duration) *WalletStore {

	return &WalletStore{db: db, queryTimeout: queryTimeout}
}

func (w *WalletStore) CreateWallet(ctx context.Context, body *CreateWalletBody) error {
	return nil
}

func (w *WalletStore) GetWallet(ctx context.Context, walletID uuid.UUID) (Wallet, error) {

	return Wallet{}, nil
}

func (w *WalletStore) TopUpWallet(ctx context.Context, walletID uuid.UUID, amount float32) error {
	return nil
}

func (w *WalletStore) ChargeWallet(ctx context.Context, walletID uuid.UUID, amount float32) error {
	return nil
}
