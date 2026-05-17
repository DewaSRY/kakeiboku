package services

import (
	"context"
	"errors"
)


var (
	ErrCreatingUserWIthDuplicateEmail = errors.New("cannot create user with duplicate email")
)

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserParams) (User, error) {
	var result User

	err := store.ExecTX(ctx, func(q *Queries) error {
		var err error

		// check if email already exist
		_, err = q.GetUserByEmail(ctx, arg.Email)
		if err == nil {
			return ErrCreatingUserWIthDuplicateEmail
		}

		result, err = q.CreateUser(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return User{}, err
	}

	return result, nil
}
