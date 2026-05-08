package test

import (
	"context"
	"math/big"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	db "github.com/dewasurya/kakeiboku/apps/apiportal/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// Helper to create a user
func createTestUser(t *testing.T, ctx context.Context) db.User {
	user, err := testStore.CreateUser(ctx, db.CreateUserParams{
		FirstName:    gofakeit.NamePrefix(),
		LastName:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		PasswordHash: "",
	})
	require.NoError(t, err)
	return user
}

// Helper to create an account
func createTestAccount(t *testing.T, ctx context.Context, userID int64, balance pgtype.Numeric, currency string) db.Account {
	acc, err := testStore.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   userID,
		Balance:  balance,
		Currency: currency,
	})
	require.NoError(t, err)
	return acc
}

func intToPgTypeNumeric(num int) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(num)), Valid: true}
}
