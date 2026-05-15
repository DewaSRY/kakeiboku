package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/mock"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransactionHandler(t *testing.T) {
	test_case := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account": gofakeit.Int64(),
				"to_account":   gofakeit.Int64(),
				"amount":       100,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(1).
					Return(&services.Transfer{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "BadRequest_MissingFields",
			body: gin.H{},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest_InvalidAmount",
			body: gin.H{
				"from_account": gofakeit.Int64(),
				"to_account":   gofakeit.Int64(),
				"amount":       -5,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError_CreateTransferTx",
			body: gin.H{
				"from_account": gofakeit.Int64(),
				"to_account":   gofakeit.Int64(),
				"amount":       100,
			},
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, errors.New("db error"))
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
			r.POST("/transaction", server.TransactionHandler)

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
