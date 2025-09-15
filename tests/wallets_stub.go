package tests

import (
	"context"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/wallet"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type StubWalletStore struct {
	Wallets []wallet.Wallet
	Fail    bool
}

func (s *StubWalletStore) CreateWallet(ctx context.Context, companyID uuid.UUID) (wallet.Wallet, error) {

	if s.Fail {
		return wallet.Wallet{}, helpers.ErrInternalServer
	}

	var walletData wallet.Wallet
	currentTime := time.Now()

	walletData.CreatedAt = &currentTime
	walletData.UpdatedAt = &currentTime
	walletData.CompanyID = companyID
	walletData.Balance = decimal.NewFromFloat(0.0)

	s.Wallets = append(s.Wallets, walletData)

	return walletData, nil
}

func (s *StubWalletStore) GetWallet(ctx context.Context, companyID uuid.UUID) (wallet.Wallet, error) {

	for _, w := range s.Wallets {
		if w.CompanyID == companyID {
			return w, nil
		}
	}

	return wallet.Wallet{}, helpers.ErrNotFound
}

func (s *StubWalletStore) TopUpWallet(ctx context.Context, companyID uuid.UUID, amount decimal.Decimal) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil

}
func (s *StubWalletStore) ChargeWallet(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal) (wallet.Wallet, error) {

	for _, w := range s.Wallets {
		if w.CompanyID == walletID {
			if w.Balance.LessThan(amount) {
				return w, helpers.ErrInsufficientFunds
			} else {
				w.Balance = w.Balance.Sub(amount)
				return w, nil
			}
		}
	}

	return wallet.Wallet{}, helpers.ErrNotFound
}
