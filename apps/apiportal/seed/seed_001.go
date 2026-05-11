package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
)

func Seed001(ctx context.Context, store services.Store) {

	// store.GetAccountByID()
	account_amount, err := store.GetAccountCount(ctx)

	if err != nil {
		log.Panic("failed to get account ammount")
	}

	if account_amount != 0 {

		fmt.Print("account already seed")
		return

	}

	create_user_list := []services.CreateUserParams{
		{
			FirstName: "Admin",
			LastName:  "Admin",
			Email:     "admin@example.com",
		},
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	user_created_list := make([]services.User, 0)

	for _, create_user := range create_user_list {
		hash_password, err := utils.HashPassword(USER_PASSWORD)
		if err != nil {
			log.Panic("failed to hash password")
		}
		create_user.PasswordHash = hash_password

		user_created, err := store.CreateUser(ctx, create_user)
		if err != nil {
			log.Panic("failed to create user")
		}

		user_created_list = append(user_created_list, user_created)
	}

	account_create_list := make([]services.Account, 0)
	for _, user_created := range user_created_list {
		create_account_params := []services.CreateAccountsParams{
			{
				UserID:   user_created.ID,
				Balance:  utils.IntToPgTypeNumeric(50),
				Currency: "USD",
			},
			{
				UserID:   user_created.ID,
				Balance:  utils.IntToPgTypeNumeric(150),
				Currency: "USD",
			},
		}
		for _, create_account_param := range create_account_params {
			account_created, err := store.CreateAccounts(ctx, create_account_param)
			if err != nil {
				log.Panic("failed to create account")
			}
			account_create_list = append(account_create_list, account_created)
		}
	}

	for _, account_created := range account_create_list {
		for _, account := range account_create_list {
			if account.ID != account_created.ID {
				_, err := store.CreateTransferTx(ctx, services.CreateTransactionParams{
					FromAccountID: account_created.ID,
					ToAccountID:   account.ID,
					Amount:        utils.IntToPgTypeNumeric(10),
				})
				if err != nil {
					log.Panic("failed to create transfer")
				}
			}
		}
	}

	log.Print("success create seed user, account, transaction ")
}
