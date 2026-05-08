package services

import (
	"context"
	"errors"
	"math/big"

	db "github.com/dewasurya/kakeiboku/apps/apiportal/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBalanceIsNegative = errors.New("account balance is negative")
)

func (t *SQLStore) CreateTransferTx(ctx context.Context, arg db.CreateTransactionParams) (*db.Transfer, error) {
	var result db.Transfer

	if err := t.ExecTX(ctx, func(q *db.Queries) error {
		// transfer money from account (sender) to account (receiver)
		fromA, err := q.GetAccountByID(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		toA, err := q.GetAccountByID(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		transferAmount := pgtype.Numeric{
			Int:              new(big.Int).Neg(arg.Amount.Int),
			Exp:              arg.Amount.Exp,
			NaN:              arg.Amount.NaN,
			InfinityModifier: arg.Amount.InfinityModifier,
			Valid:            arg.Amount.Valid,
		}

		if fromA.ID < toA.ID {
			_, err := _addBalance(ctx, q, _AddBalanceArgs{
				accountId: fromA.ID,
				amount:    transferAmount,
			})
			if err != nil {
				return err
			}

			_, err = _addBalance(ctx, q, _AddBalanceArgs{
				accountId: toA.ID,
				amount:    arg.Amount,
			})
			if err != nil {
				return err
			}

		} else {
			_, err := _addBalance(ctx, q, _AddBalanceArgs{
				accountId: toA.ID,
				amount:    arg.Amount,
			})
			if err != nil {
				return err
			}

			_, err = _addBalance(ctx, q, _AddBalanceArgs{
				accountId: fromA.ID,
				amount:    transferAmount,
			})
			if err != nil {
				return err
			}
		}

		// create the transaction
		result, err = q.CreateTransaction(ctx, db.CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &result, nil

}

type _AddBalanceArgs struct {
	accountId int64
	amount    pgtype.Numeric
}

func _addBalance(ctx context.Context, q *db.Queries, arg _AddBalanceArgs) (account db.Account, err error) {
	account, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:     arg.accountId,
		Amount: arg.amount,
	})

	if account.Balance.Int.Cmp(big.NewInt(0)) < 0 {
		err = ErrBalanceIsNegative
		return
	}

	return
}
