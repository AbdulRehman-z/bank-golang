package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

			server := NewServer(store)

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
				Owner:    account.Owner,
				Currency: "invalid currency",
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
				fmt.Printf("b")
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
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)

				payload, err := json.Marshal(tc.body)
				require.NoError(t, err)

				request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(payload))
				request.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)

				response, err := server.router.Test(request)
				require.NoError(t, err)
				tc.checkResp(t, response)
			})
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name       string
		query      any
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "ListAccountsSuccuess",
			query: types.ListAccountsRequest{
				PageID:   1,
				PageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{account}, nil)
			},
			checkResp: func(t *testing.T, resp *http.Response) {

				require.Equal(t, http.StatusOK, resp.StatusCode)
				// requireBodyMatchAccount(t, resp.Body, account)
				requireBodyMatch(t, resp.Body, []db.Account{account})
			},
		},
		// {
		// 	name: "InvalidQuery",
		// },
		// {
		// 	name: "DatabaseError",
		// },
		// {
		// 	name: "NoAccountsFound",
		// },
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)

			url := fmt.Sprintf("/accounts")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			response, err := server.router.Test(request)

			tc.checkResp(t, response)
		})
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(util.GenerateRandomInt(1, 1000)),
		Owner:    util.GenerateRandomString(10),
		Balance:  int64(util.GenerateRandomMoney()),
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
		require.Equal(t, expected.Balance, int64(account["balance"].(float64)))
	case []db.Account:
		accounts, ok := response.Data.([]interface{})
		require.True(t, ok)
		require.Equal(t, len(expected), len(accounts))
		for i, account := range expected {
			accountData, ok := accounts[i].(map[string]interface{})
			require.True(t, ok)
			require.Equal(t, account.ID, int64(accountData["id"].(float64)))
			require.Equal(t, account.Owner, accountData["owner"].(string))
			require.Equal(t, account.Balance, int64(accountData["balance"].(float64)))
		}
	default:
		t.Fatalf("unexpected type: %T", expected)
	}
}

// func TestListAccountsHandler(t *testing.T) {
//     account := randomAccount()

//     testCases := []struct {
//         name       string
//         query      any
//         buildStubs func(store *mockdb.MockStore)
//         checkResp  func(t *testing.T, resp *http.Response)
//     }{
//         {
//             name: "ListAccounts_Success",
//             query: types.ListAccountsRequest{
//                 PageID:   1,
//                 PageSize: 5,
//             },
//             buildStubs: func(store *mockdb.MockStore) {
//                 store.EXPECT().
//                     ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{account}, nil)
//             },
//             checkResp: func(t *testing.T, resp *http.Response) {
//                 require.Equal(t, http.StatusOK, resp.StatusCode)
//                 requireBodyMatchAccount(t, resp.Body, account)
//             },
//         },
//         {
//             name: "ListAccounts_InvalidQuery",
//         },
//         {
//             name: "ListAccounts_DatabaseError",
//         },
//         {
//             name: "ListAccounts_NoAccountsFound",
//         },
//     }

//     for _, tc := range testCases {
//         t.Run(tc.name, func(t *testing.T) {
//             ctrl := gomock.NewController(t)
//             defer ctrl.Finish()

//             store := mockdb.NewMockStore(ctrl)
//             tc.buildStubs(store)

//             server := NewServer(store)

//             req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
//             q := req.URL.Query()
//             for key, val := range tc.query.(map[string][]string) {
//                 q.Add(key, val[0])
//             }
//             req.URL.RawQuery = q.Encode()

//             resp := httptest.NewRecorder()

//             server.listAccountsHandler(resp, req)

//             tc.checkResp(t, resp.Result())
//         })
//     }
// }
