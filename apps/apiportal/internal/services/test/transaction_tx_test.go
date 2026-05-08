package test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	db "github.com/dewasurya/kakeiboku/apps/apiportal/internal/database/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferTx_Integration(t *testing.T) {
	ctx := context.Background()

	user, err := testStore.CreateUser(ctx, db.CreateUserParams{
		FirstName:    gofakeit.NamePrefix(),
		LastName:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		PasswordHash: "",
	})
	require.NoError(t, err)

	// Create two accounts
	fromAccount, err := testStore.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   user.ID,
		Balance:  intToPgTypeNumeric(100),
		Currency: "USD",
	})

	require.NoError(t, err)

	toAccount, err := testStore.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   user.ID,
		Balance:  intToPgTypeNumeric(50),
		Currency: "USD",
	})

	require.NoError(t, err)

	amount := intToPgTypeNumeric(10)

	// test for concurrent situation
	const num = 5
	errors := make(chan error, num)
	result := make(chan *db.Transfer, num)

	for i := 0; i < num; i++ {
		go func() {
			transfer, err := testStore.CreateTransferTx(ctx, db.CreateTransactionParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errors <- err
			result <- transfer
		}()
	}

	for i := 0; i < num; i++ {
		err := <-errors
		require.NoError(t, err)

		transfer := <-result
		require.NoError(t, err)
		require.NotNil(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
	}

}

func TestTransferTx_InsufficientBalance_Integration(t *testing.T) {
	ctx := context.Background()
	user, err := testStore.CreateUser(ctx, db.CreateUserParams{
		FirstName:    gofakeit.NamePrefix(),
		LastName:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		PasswordHash: "",
	})
	require.NoError(t, err)

	// Create two accounts
	fromAccount, err := testStore.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   user.ID,
		Balance:  intToPgTypeNumeric(100),
		Currency: "USD",
	})

	require.NoError(t, err)

	toAccount, err := testStore.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   user.ID,
		Balance:  intToPgTypeNumeric(50),
		Currency: "USD",
	})

	require.NoError(t, err)

	amount := intToPgTypeNumeric(1000)

	transfer, err := testStore.CreateTransferTx(ctx, db.CreateTransactionParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	})
	require.Nil(t, transfer)
	require.Error(t, err)

}

