package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user *User) error
}

// CreateUserTxResult is the result of the create user tracsaction
type CreateUserTxResult struct {
	User *User
}

// CreateUserTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts'
// balance within a single database transaction
func (s *SqlStore) CreateUserTx(ctx context.Context, arg *CreateUserTxParams) (*CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.ExecTx(ctx, func(q *Queries) error {

		user, err := q.CreateUser(ctx, &arg.CreateUserParams)
		if err != nil {
			return err
		}

		err = arg.AfterCreate(user)
		if err != nil {
			return err
		}

		result.User = user
		return nil
	})

	return &result, err
}
