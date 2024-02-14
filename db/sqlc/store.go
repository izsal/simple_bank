package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters for the TransferTx function
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount 		  int64 `json:"amount"`
}

// TransferTxResult is the result of the TransferTx function
type TransferTxResult struct {
	Transfer 	Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount 	Account `json:"to_account"`
	FromEntry 	Entry `json:"from_entry"`
	ToEntry 	Entry `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// get account -> update its balance
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}


		return nil
		})
		return result, err
	}
