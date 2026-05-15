package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/mock"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	test_case := []struct {
		name          string
		buildStubs    func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK_StatusUp",
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().Health(gomock.Any()).Times(1).Return(services.StoreHealthRecord{
					Status:  "up",
					Message: "It's healthy",
				})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "ServiceUnavailable_StatusDown",
			buildStubs: func(store *mock.MockStore, tokenMaker *mock.MockTokenMaker) {
				store.EXPECT().Health(gomock.Any()).Times(1).Return(services.StoreHealthRecord{
					Status: "down",
					Error:  "connection refused",
				})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
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
			r.GET("/health", server.HealthHandler)

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
