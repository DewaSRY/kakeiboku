package test

import (
	"context"
	"testing"

	db "github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/stretchr/testify/require"
)



func TestCreateTransferTx_Integration(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, ctx)
	fromAccount := createTestAccount(t, ctx, user.ID, intToPgTypeNumeric(100), "USD")
	toAccount := createTestAccount(t, ctx, user.ID, intToPgTypeNumeric(50), "USD")
	amount := intToPgTypeNumeric(10)

	const num = 5
	errors := make(chan error, num)
	results := make(chan *db.Transfer, num)

	for i := 0; i < num; i++ {
		go func() {
			transfer, err := testStore.CreateTransferTx(ctx, db.CreateTransactionParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errors <- err
			results <- transfer
		}()
	}

	for i := 0; i < num; i++ {
		err := <-errors
		require.NoError(t, err)

		transfer := <-results
		require.NotNil(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
	}
}

func TestTransferTx_InsufficientBalance_Integration(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, ctx)
	fromAccount := createTestAccount(t, ctx, user.ID, intToPgTypeNumeric(100), "USD")
	toAccount := createTestAccount(t, ctx, user.ID, intToPgTypeNumeric(50), "USD")
	amount := intToPgTypeNumeric(1000)

	transfer, err := testStore.CreateTransferTx(ctx, db.CreateTransactionParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	})
	require.Nil(t, transfer)
	require.Error(t, err)
}

