package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

type Store interface {
	CreateWallet(ctx context.Context, companyID uuid.UUID) (Wallet, error)
	GetWallet(ctx context.Context, userID uuid.UUID) (Wallet, error)
	TopUpWallet(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (Wallet, error)
	ChargeWallet(ctx context.Context, companyID uuid.UUID, amount decimal.Decimal) (Wallet, error)
}

const UniqueViolationCode = "23505"

type WalletStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewWalletStore(db *pgxpool.Pool, queryTimeout time.Duration) *WalletStore {

	return &WalletStore{db: db, queryTimeout: queryTimeout}
}

func (w *WalletStore) CreateWallet(ctx context.Context, companyID uuid.UUID) (Wallet, error) {
	ctx, cancel := w.WithTimeout(ctx)
	defer cancel()

	query := `INSERT INTO wallets (balance, company_id, updated_at) 
				VALUES (@balance, @company_id, NOW()) 
				RETURNING id, balance, company_id, created_at, updated_at;`
	args := pgx.NamedArgs{
		"balance":    decimal.Zero,
		"company_id": companyID,
	}

	var wallet Wallet

	row := w.db.QueryRow(ctx, query, args)

	err := row.Scan(&wallet.ID, &wallet.Balance, &wallet.CompanyID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return Wallet{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Wallet{}, fmt.Errorf("error scanning row (create wallet): %w", err)
	}

	return wallet, nil
}

func (w *WalletStore) GetWallet(ctx context.Context, userID uuid.UUID) (Wallet, error) {
	ctx, cancel := w.WithTimeout(ctx)
	defer cancel()

	query := `SELECT w.id, w.balance, w.company_id, w.created_at, w.updated_at 
				FROM companies c 
				JOIN wallets w ON w.company_id = c.id
				WHERE user_id = @userID;`
	args := pgx.NamedArgs{
		"userID": userID,
	}

	var wallet Wallet

	row := w.db.QueryRow(ctx, query, args)

	err := row.Scan(&wallet.ID, &wallet.Balance, &wallet.CompanyID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Wallet{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Wallet{}, fmt.Errorf("error scanning row (get wallet): %w", err)
	}

	return wallet, nil
}

func (w *WalletStore) TopUpWallet(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (Wallet, error) {
	ctx, cancel := w.WithTimeout(ctx)
	defer cancel()

	query := `UPDATE wallets w
				SET balance = w.balance + @amount
				FROM companies c
				WHERE w.company_id = c.id
  				AND c.user_id = @userID
				RETURNING w.id, w.balance, w.company_id, w.created_at, w.updated_at;
			`
	args := pgx.NamedArgs{
		"amount": amount,
		"userID": userID,
	}

	var wallet Wallet

	row := w.db.QueryRow(ctx, query, args)

	err := row.Scan(&wallet.ID, &wallet.Balance, &wallet.CompanyID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		err = errors.Join(helpers.ErrInternalServer, err)
		return Wallet{}, fmt.Errorf("error scanning row (topup wallet): %w", err)
	}

	return wallet, nil
}

func (w *WalletStore) ChargeWallet(ctx context.Context, companyID uuid.UUID, amount decimal.Decimal) (Wallet, error) {
	ctx, cancel := w.WithTimeout(ctx)
	defer cancel()

	query := `UPDATE wallets 
				SET balance = balance - @amount 
				WHERE company_id = @companyID AND balance >= @amount
				RETURNING id, balance, company_id, created_at, updated_at;`
	args := pgx.NamedArgs{
		"amount":    amount,
		"companyID": companyID,
	}

	var wallet Wallet

	row := w.db.QueryRow(ctx, query, args)

	err := row.Scan(&wallet.ID, &wallet.Balance, &wallet.CompanyID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Wallet{}, helpers.ErrInsufficientFunds
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Wallet{}, fmt.Errorf("error scanning row (charge wallet): %w", err)
	}

	return wallet, nil
}

func (w *WalletStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, w.queryTimeout)
}
