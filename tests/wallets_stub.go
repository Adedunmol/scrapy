package tests

import (
	"context"
	"github.com/Adedunmol/scrapy/api/wallet"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StubWalletStore struct{}

func (s *StubWalletStore) CreateWallet(ctx context.Context, companyID uuid.UUID) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil
}

func (s *StubWalletStore) GetWallet(ctx context.Context, companyID uuid.UUID) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil
}

func (s *StubWalletStore) TopUpWallet(ctx context.Context, companyID uuid.UUID, amount decimal.Decimal) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil

}
func (s *StubWalletStore) ChargeWallet(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil
}
