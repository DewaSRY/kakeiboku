package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/server"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)



func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func newTestServer(store services.Store, tokenMaker token.TokenMaker) *server.Server {
	config := utils.Config{
		SecretKey:           gofakeit.UUID(),
		AccessTokenDuration: time.Minute,
	}

	mock_server := &server.Server{
		Port:  4000,
		Store: store,
		Config: config,
		Token: tokenMaker,
	}

	server.RegisterRoutes(mock_server)
	return mock_server
}


type AddAuthorizationArg struct {
	UserId    int64
	Email     string
	TokenType token.TokenType
	Duration  time.Duration
}



func AddAuthorization(t *testing.T, request *http.Request, tokenMaker token.TokenMaker, arg AddAuthorizationArg) {
	token, payload, err := tokenMaker.CreateToken(
		arg.UserId,
		arg.Email,
		arg.Duration,
		arg.TokenType,
	)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	authorizationHeader := fmt.Sprintf("Bearer %s", token)
	request.Header.Set("Authorization", authorizationHeader)

}