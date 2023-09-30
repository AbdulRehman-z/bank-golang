package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	mockdb "github.com/AbdulRehman-z/bank-golang/db/mock"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {

	account := randomAccount()

	testcases := []struct {
		name       string
		accountID  int64
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name:      "GetAccountSuccess",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				// requireBodyMatchAccount(t, response.Body, account)
				requireBodyMatch(t, response.Body, account)
			},
		},
		{
			name:      "AccountNotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:      "InvalidRequestField",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:      "DatabaseError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			response, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResp(t, response)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name       string
		body       any
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "CreateAccountSuccess",
			body: types.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)
				// requireBodyMatchAccount(t, resp.Body, account)
				requireBodyMatch(t, resp.Body, account)
			},
		},
		{
			name: "InvalidRequestFields",
			body: types.CreateAccountRequest{
				Owner:    "jj",
				Currency: account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, resp.StatusCode, http.StatusBadRequest)
			},
		},
		{
			name: "InvalidRequestBody",
			body: "DAWDW",

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, resp.StatusCode, http.StatusBadRequest)
			},
		},
		{
			name: "DatabaseError",
			body: types.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			payload, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			response, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResp(t, response)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name       string
		query      interface{}
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "ListAccountsSuccuess",
			query: map[string]string{
				"page_Id":   "2",
				"page_Size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{account}, nil)
			},
			checkResp: func(t *testing.T, resp *http.Response) {

				require.Equal(t, http.StatusOK, resp.StatusCode)
				requireBodyMatch(t, resp.Body, []db.Account{account})
			},
		},
		{
			name: "InvalidQueryValue",
			query: map[string]string{
				"page_id":   "0",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:  "InvalidQueryBody",
			query: map[string]string{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name: "DatabaseError",
			query: map[string]string{
				"page_id":   "20000",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
		{
			name: "NoAccountsFound",
			query: map[string]string{
				"page_id":   "20202",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			params := url.Values{}
			for k, v := range tc.query.(map[string]string) {
				params.Add(k, v)
			}

			encodedUrl := fmt.Sprintf("/accounts?%s", params.Encode())
			request, err := http.NewRequest(http.MethodGet, encodedUrl, nil)
			require.NoError(t, err)

			response, err := server.router.Test(request)
			require.NoError(t, err)

			tc.checkResp(t, response)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	account := randomAccount()

	// convert balance(string) to balance(string)
	balance, err := strconv.ParseFloat(account.Balance, 64)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		body       any
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "UpdateAccountSuccess",
			body: types.UpdateAccountRequest{
				ID:      int64(2),
				Balance: int64(balance),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				// requireBodyMatchAccount(t, resp.Body, account)
				requireBodyMatch(t, resp.Body, account)
			},
		},
		{
			name: "InvalidRequestBody",
			body: "DAWDW",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name: "InvalidRequestFields",
			body: types.UpdateAccountRequest{
				ID:      int64(2),
				Balance: -1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name: "DatabaseError",
			body: types.UpdateAccountRequest{
				ID:      int64(2),
				Balance: int64(balance),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			payload, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/accounts", bytes.NewBuffer(payload))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			response, err := server.router.Test(request)
			fmt.Printf("response: %v\n", response)
			fmt.Printf("err: %v\n", err)
			require.NoError(t, err)

			tc.checkResp(t, response)
		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {

	testCases := []struct {
		name       string
		accountID  int64
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name:      "DeleteAccountSuccess",
			accountID: 2,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:      "InvalidRequestField",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:      "DatabaseError",
			accountID: 2,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			response, err := server.router.Test(request)
			require.NoError(t, err)

			tc.checkResp(t, response)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.GenerateRandomInt(1, 1000),
		Owner:    util.GenerateRandomString(10),
		Balance:  util.GenerateRandomMoney(),
		Currency: util.GenerateRandomCurrencyCode(),
	}
}

func requireBodyMatch(t *testing.T, body io.Reader, expected interface{}) {
	var response struct {
		Data interface{} `json:"data"`
	}
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)

	switch expected := expected.(type) {
	case db.Account:
		account, ok := response.Data.(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, expected.ID, int64(account["id"].(float64)))
		require.Equal(t, expected.Owner, account["owner"].(string))
		require.Equal(t, expected.Balance, account["balance"].(string))
	case []db.Account:
		accounts, ok := response.Data.([]interface{})
		require.True(t, ok)
		// require.Equal(t, len(expected), len(accounts))
		for i, account := range expected {
			accountData, ok := accounts[i].(map[string]interface{})
			require.True(t, ok)
			require.Equal(t, account.ID, int64(accountData["id"].(float64)))
			require.Equal(t, account.Owner, accountData["owner"].(string))
			require.Equal(t, account.Balance, accountData["balance"].(string))
		}
	default:
		t.Fatalf("unexpected type: %T", expected)
	}
}
