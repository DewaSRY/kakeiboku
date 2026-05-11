package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/middleware"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/mock"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)



func TestAuthMiddleWare(t *testing.T) {

	test_case := []struct {
		name          string
		setupAuth     func(request *http.Request, tokenMaker token.TokenMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request, tokenMaker token.TokenMaker) {
				AddAuthorization(t, request, tokenMaker, AddAuthorizationArg{
					UserId:    1,
					Email:     gofakeit.Email(),
					TokenType: token.TokenTypeAccessToken,
					Duration:  time.Minute,
				})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotAuthorization",
			setupAuth: func(request *http.Request, tokenMaker token.TokenMaker) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request, tokenMaker token.TokenMaker) {
				AddAuthorization(t, request, tokenMaker, AddAuthorizationArg{
					UserId:    1,
					Email:     gofakeit.Email(),
					TokenType: token.TokenTypeAccessToken,
					Duration:  -time.Minute,
				})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range test_case {
		tc := test_case[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			token_maker, err := token.NewJWTMaker(gofakeit.UUID())
			require.NoError(t, err)
			server := newTestServer(mock.NewMockStore(ctrl), token_maker)

			auth_path := "/test/auth"

			// Create a new gin router
			router := gin.New()
			router.GET(auth_path, middleware.AuthMiddleware(server.Token), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create a new HTTP request
			request, err := http.NewRequest(http.MethodGet, auth_path, nil)
			require.NoError(t, err)

			// Setup authorization
			tc.setupAuth(request, server.Token)

			// Create a response recorder
			recorder := httptest.NewRecorder()
			// Call the handler
			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
