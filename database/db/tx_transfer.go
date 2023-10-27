package db

import "context"

// TransferTxResult is the result of the transfer tracsaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `jsno:"to_account"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts'
// balance within a single database transaction
func (s *SqlStore) TransferTx(ctx context.Context, arg *CreateTransferParams) (*TransferTxResult, error) {
	var result TransferTxResult

	err := s.ExecTx(ctx, func(q *Queries) error {

		transfer, err := q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}

		fromEntry, err := q.CreateEntry(ctx, &CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		toEntry, err := q.CreateEntry(ctx, &CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		var fromAccount *Account
		var toAccount *Account

		if arg.FromAccountID < arg.ToAccountID {
			fromAccount, toAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			toAccount, fromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		result.Transfer = *transfer
		result.FromEntry = *fromEntry
		result.ToEntry = *toEntry
		result.FromAccount = *fromAccount
		result.ToAccount = *toAccount

		return nil
	})

	return &result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1 *Account, account2 *Account, err error) {
	account1, err = q.AddAccountBalance(ctx, &AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, &AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
