package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/mock"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newFakePayload(email string) *token.Payload {
	payload, _ := token.NewPayload(1, email, time.Minute, token.TokenTypeAccessToken)
	return payload
}

func TestSignUpHandler(t *testing.T) {
	test_case := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker)

		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": gofakeit.Password(true, true, true, true, false, 12),
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				var user services.User
				gofakeit.Struct(&user)
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2).Return("token", newFakePayload(user.Email), nil)
				store.EXPECT().SetSession(gomock.Any(), gomock.Any()).Times(1).Return(services.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest_MissingFields",
			body: gin.H{},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError_CreateUser",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": gofakeit.Password(true, true, true, true, false, 12),
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).
					Return(services.User{}, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Conflict_DuplicateEmail",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": gofakeit.Password(true, true, true, true, false, 12),
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).
					Return(services.User{}, services.ErrCreatingUserWIthDuplicateEmail)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InternalError_CreateToken",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": gofakeit.Password(true, true, true, true, false, 12),
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				var user services.User
				gofakeit.Struct(&user)
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("", nil, errors.New("token error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError_SetSession",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": gofakeit.Password(true, true, true, true, false, 12),
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				var user services.User
				gofakeit.Struct(&user)
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2).Return("token", newFakePayload(user.Email), nil)
				store.EXPECT().SetSession(gomock.Any(), gomock.Any()).Times(1).
					Return(services.Session{}, errors.New("session error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range test_case {
		tc := test_case[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tokenMaker := mock.NewMockTokenMaker(ctrl)
			tc.buildStubs(store, tokenMaker)

			server := newTestServer(store, tokenMaker)
			r := gin.New()
			r.POST("/signup", server.SignUpHandler)

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	password := "Password123!"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	fakeUser := services.User{
		ID:           gofakeit.Int64(),
		Email:        gofakeit.Email(),
		PasswordHash: hashedPassword,
	}

	test_case := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    fakeUser.Email,
				"password": password,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeUser.Email).Times(1).Return(fakeUser, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2).Return("token", newFakePayload(fakeUser.Email), nil)
				store.EXPECT().SetSession(gomock.Any(), gomock.Any()).Times(1).Return(services.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest_MissingFields",
			body: gin.H{},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized_UserNotFound",
			body: gin.H{
				"email":    gofakeit.Email(),
				"password": password,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).
					Return(services.User{}, errors.New("user not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized_WrongPassword",
			body: gin.H{
				"email":    fakeUser.Email,
				"password": "wrongpassword",
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeUser.Email).Times(1).Return(fakeUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError_CreateToken",
			body: gin.H{
				"email":    fakeUser.Email,
				"password": password,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeUser.Email).Times(1).Return(fakeUser, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("", nil, errors.New("token error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError_SetSession",
			body: gin.H{
				"email":    fakeUser.Email,
				"password": password,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeUser.Email).Times(1).Return(fakeUser, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2).Return("token", newFakePayload(fakeUser.Email), nil)
				store.EXPECT().SetSession(gomock.Any(), gomock.Any()).Times(1).
					Return(services.Session{}, errors.New("session error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range test_case {
		tc := test_case[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tokenMaker := mock.NewMockTokenMaker(ctrl)
			tc.buildStubs(store, tokenMaker)

			server := newTestServer(store, tokenMaker)
			r := gin.New()
			r.POST("/login", server.LoginHandler)

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestRefreshTokenHandler(t *testing.T) {
	fakeEmail := gofakeit.Email()
	fakePayload := newFakePayload(fakeEmail)

	fakeUser := services.User{
		ID:    gofakeit.Int64(),
		Email: fakeEmail,
	}

	test_case := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"refresh_token": "sometoken",
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				tokenMaker.EXPECT().VerifyToken("sometoken", token.TokenTypeRefreshToken).Times(1).Return(fakePayload, nil)
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeEmail).Times(1).Return(fakeUser, nil)
				tokenMaker.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("new_access_token", newFakePayload(fakeEmail), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest_MissingToken",
			body: gin.H{},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				tokenMaker.EXPECT().VerifyToken(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized_InvalidToken",
			body: gin.H{
				"refresh_token": "invalidtoken",
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				tokenMaker.EXPECT().VerifyToken("invalidtoken", token.TokenTypeRefreshToken).Times(1).
					Return(nil, token.ErrInvalidToken)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized_UserNotFound",
			body: gin.H{
				"refresh_token": "sometoken",
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				tokenMaker.EXPECT().VerifyToken("sometoken", token.TokenTypeRefreshToken).Times(1).Return(fakePayload, nil)
				store.EXPECT().GetUserByEmail(gomock.Any(), fakeEmail).Times(1).
					Return(services.User{}, errors.New("user not found"))
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

			store := mock.NewMockStore(ctrl)
			tokenMaker := mock.NewMockTokenMaker(ctrl)
			tc.buildStubs(store, tokenMaker)

			server := newTestServer(store, tokenMaker)
			r := gin.New()
			r.POST("/refresh", server.RefreshTokenHandler)

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
