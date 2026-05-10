package tests

// func TestSingUpHandler(t *testing.T) {

// 	test_case := []struct {
// 		name string
// 		body gin.H
// 		buildDbStubs func(store *mock.MockStore)
// 		buildTokenStubs func(token *mock.MockTokenMaker)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"email": gofakeit.Email(),
// 				"password": gofakeit.Password(true, true, true, true, false, 12),
// 			},
// 			buildDbStubs: func(store *mock.MockStore) {
// 				var user_created services.User
// 				gofakeit.Struct(&user_created)

// 				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).
// 				Return(user_created, nil)

// 			},
// 			buildTokenStubs: func(token *mock.MockTokenMaker) {

// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}

// 		for i := range test_case {
// 			tc := test_case[i]

// 			t.Run(tc.name, func(t *testing.T) {
// 				ctrl := gomock.NewController(t)
// 				defer ctrl.Finish()

// 				store := mock.NewMockStore(ctrl)
// 				tc.buildDbStubs(store)
// 				token_maker := mock.NewMockTokenMaker(ctrl)

// 				server := newTestServer(store, token_maker)
// 				r := gin.Default()
// 				r.POST("/signup", server.SignUpHandler)

// 				reqBody, err := json.Marshal(tc.body)
// 				require.NoError(t, err)

// 				req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
// 				require.NoError(t, err)
// 				req.Header.Set("Content-Type", "application/json")

// 				recorder := httptest.NewRecorder()
// 				r.ServeHTTP(recorder, req)

// 				tc.checkResponse(recorder)
// 			})
// 		}

// }