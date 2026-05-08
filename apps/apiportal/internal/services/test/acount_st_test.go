package test

import (
	"context"
	"testing"

	db "github.com/dewasurya/kakeiboku/apps/apiportal/internal/database/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T){
	ctx:= context.Background() 

	user := createTestUser(t,ctx)



	arg:=  db.CreateAccountsParams{
		UserID: user.ID,
		Balance: intToPgTypeNumeric(0),
		Currency: "",
	}

	acc, err := testStore.CreateAccounts(ctx, arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, arg.UserID, acc.UserID)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)
}