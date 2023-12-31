package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "tutorial.sqlc.dev/app/db/mock"
	db "tutorial.sqlc.dev/app/db/sqlc"
	mytoken "tutorial.sqlc.dev/app/token"
	"tutorial.sqlc.dev/app/util"
)

func TestGetAccount(t *testing.T) {
	randomUser, _ := randomUser(t)
	account := randomAccount(randomUser.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyAccountMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "Not Found",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			accountID: 0,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:      "Forbidden",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "randomuser", time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() //important to check whether all methods that were expected to be called were actually called.
			store := mockdb.NewMockStore(ctrl)
			
			//build stubs
			tc.buildStubs(store)

			//start a test server
			server, err := NewTestServer(t, store)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/accounts/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setUpAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req) //this would send `req` through the server router and record its response in `recorder`.

			//check response
			tc.checkResponse(t, recorder)

		})

	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
}

func RequireBodyAccountMatch(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
