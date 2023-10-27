package db

import "context"

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        *User
	VerifyEmail *VerifyEmail
}

func (s *SqlStore) VerifyEmailTx(ctx context.Context, arg *VerifyEmailTxParams) (*VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := s.ExecTx(ctx, func(q *Queries) error {

		verifyEmail, err := q.UpdateVerifyEmail(ctx, &UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		isEmailVerified := true

		user, err := q.UpdateUser(ctx, &UpdateUserParams{
			Username:        verifyEmail.Username,
			IsEmailVerified: &isEmailVerified,
		})
		if err != nil {
			return err
		}

		result.User = user
		result.VerifyEmail = verifyEmail
		return nil
	})

	return &result, err
}
