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

func TestCreateUserAPI(t *testing.T) {

	user := randomUser()

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "CreateUserSuccess",
			body: types.CreateUserRequest{
				Username:       util.GenerateRandomOwnerName(),
				FullName:       util.GenerateRandomOwnerName(),
				Email:          util.GenerateRandomEmail(),
				HashedPassword: "secretrrr",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusCreated, response.StatusCode)
				requireUserBodyMatch(t, response.Body, user)
			},
		},
		{
			name: "InvalidRequestBody",
			body: "akd",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidRequestFields",
			body: types.CreateUserRequest{
				Username:       util.GenerateRandomOwnerName(),
				FullName:       util.GenerateRandomOwnerName(),
				Email:          util.GenerateRandomEmail(),
				HashedPassword: "se",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			body: types.CreateUserRequest{
				Username:       util.GenerateRandomOwnerName(),
				FullName:       util.GenerateRandomOwnerName(),
				Email:          util.GenerateRandomEmail(),
				HashedPassword: "secretrrr",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "UserAlreadyExists",
			body: types.CreateUserRequest{
				Username:       util.GenerateRandomOwnerName(),
				FullName:       util.GenerateRandomOwnerName(),
				Email:          util.GenerateRandomEmail(),
				HashedPassword: "secretrrr",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		server := NewServer(store)
		tc.buildStubs(store)

		body, err := json.Marshal(tc.body)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		response, err := server.router.Test(request)
		require.NoError(t, err)

		tc.checkResponse(t, response)

	}
}

func TestGetUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "GetUserSuccess",
			body: "akd",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireUserBodyMatch(t, response.Body, user)
			},
		},
		{
			name: "InternalError",
			body: "akd",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "UserNotFound",
			body: "akd",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name: "InvalidParams",
			body: "a",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		server := NewServer(store)
		tc.buildStubs(store)

		// get username parameter from body
		username := tc.body.(string)

		url := "/users/" + username

		request, err := http.NewRequest(http.MethodGet, url, nil)
		fmt.Printf("url: %v\n", url)
		require.NoError(t, err)

		response, err := server.router.Test(request)
		require.NoError(t, err)

		tc.checkResponse(t, response)
	}
}

func randomUser() db.User {
	return db.User{
		Username:       util.GenerateRandomOwnerName(),
		FullName:       util.GenerateRandomOwnerName(),
		Email:          util.GenerateRandomEmail(),
		HashedPassword: "secreterr",
	}
}

func requireUserBodyMatch(t *testing.T, body io.Reader, expected interface{}) {
	var userResponse struct {
		Data db.User `json:"data"`
	}
	err := json.NewDecoder(body).Decode(&userResponse)
	require.NoError(t, err)

	require.Equal(t, expected, userResponse.Data)
}
