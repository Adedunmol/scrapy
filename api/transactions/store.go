package transactions

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	CreateTransaction(ctx context.Context, body *CreateTransactionBody) (Transaction, error)
}

type TransactionStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewTransactionStore(db *pgxpool.Pool, queryTimeout time.Duration) *TransactionStore {

	return &TransactionStore{db: db, queryTimeout: queryTimeout}
}

func (t *TransactionStore) CreateTransaction(ctx context.Context, body *CreateTransactionBody) (Transaction, error) {
	ctx, cancel := t.WithTimeout(ctx)
	defer cancel()

	query := `INSERT INTO transactions (amount, balance_before, balance_after, reference, status, type, wallet_id) 
				VALUES (@amount, @balanceBefore, @balanceAfter, @reference, @status, @type, @walletID) 
				RETURNING id, amount, balance_before, balance_after, reference, status, type, wallet_id, created_at;`
	args := pgx.NamedArgs{
		"amount":        body.Amount,
		"balanceBefore": body.BalanceBefore,
		"balanceAfter":  body.BalanceAfter,
		"reference":     body.Reference,
		"status":        body.Status,
		"walletID":      body.WalletID,
		"type":          body.Type,
	}

	var transaction Transaction

	row := t.db.QueryRow(ctx, query, args)

	err := row.Scan(&transaction.ID, &transaction.Amount, &transaction.BalanceBefore, &transaction.BalanceAfter, &transaction.Reference, &transaction.Status, &transaction.Type, &transaction.WalletID, &transaction.CreatedAt)

	if err != nil {
		err = errors.Join(helpers.ErrInternalServer, err)
		return Transaction{}, fmt.Errorf("error scanning row (create transaction): %w", err)
	}

	return Transaction{}, nil
}

func (t *TransactionStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, t.queryTimeout)
}
