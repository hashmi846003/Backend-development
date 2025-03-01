/*package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}
func NewStore(db *sql.DB) *Store{
	return &Store{
	db:	db,
	Queries:	New(db),
}

}
func (store *Store)execTx(ctx context.Context,fn func(*Queries)error)error{
	tx,err:=store.db.BeginTx(ctx,nil)
	if err!=nil{
		return err
	}
	q:=New(tx)
	err=fn(q)
	if err!=nil{
		if rbErr:=tx.Rollback();rbErr!=nil{
			return fmt.Errorf("tx err:%v,rb err:%v",err,rbErr)
		}
		return err
	}
	return tx.Commit()
}
type TransferTxParams struct{
	FromAccountID 		int64 		`json:"from_account_id"`
	ToAccountID 		int64 		`json:"to_account_id"`
	Amount 				int64		`json:"amount"`
}
type TransferTxResult  struct{
Transfer	Transfer	`json:"Transfer"`
FromAccount	Account		`json:"from_account"`
ToAccount	Account     `json:"to_account"`
FromEntry	Entry		`json:"from_entry"`
ToEntry     Entry       `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult,error){
var result	TransferTxResult

err:=store.execTx(ctx,func(q *Queries)error{
var err error
result.Transfer,err=q.CreateTransfer(ctx,CreateTransferParams{
FromAccountID:	arg.FromAccountID,
ToAccountID:	arg.ToAccountID,
Amount:			arg.Amount,
})
if err!=nil{
return err
}
result.FromEntry,err=q.CreateEntry(ctx,CreateEntryParams{
AccountID: arg.FromAccountID,
Amount:	-arg.Amount,
})
if err!=nil{
return err
}
result.ToEntry,err=q.CreateEntry(ctx,CreateEntryParams{
AccountID: arg.ToAccountID,
Amount:	arg.Amount,
})
if err!=nil{
return err
}
account1,err:=q.GetAccountForUpdate(ctx,arg.FromAccountID)
if err!=nil{
	return err
}
result.FromAccount,err:=q.UpdateAccount(ctx,UpdateAccountParams{
	ID: arg.FromAccountID,
	Balance: account1.Balance-arg.Amount,
})
if err!=nil{
	return err
}
account2,err:=q.GetAccountForUpdate(ctx,arg.ToAccountID)
if err!=nil{
	return err
}
result.ToAccount,err:=q.UpdateAccount(ctx,UpdateAccountParams{
	ID: arg.ToAccountID,
	Balance: account2.Balance+arg.Amount,
})
if err!=nil{
	return err

return nil
})
return result,err
}*/
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store struct wraps the database and queries
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore initializes a new Store instance
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams defines parameters for a money transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult defines the result of a money transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer within a transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create debit entry for sender account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create credit entry for receiver account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Fetch sender account with row-level lock
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		// Update sender account balance
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// Fetch receiver account with row-level lock
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		// Update receiver account balance
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		// Assign updated accounts to result
		result.FromAccount = account1
		result.ToAccount = account2

		return nil
	}) // Ensuring execTx function is properly closed

	return result, err
} // Ensuring TransferTx function is properly closed
