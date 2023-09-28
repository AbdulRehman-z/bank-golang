package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/AbdulRehman-z/bank-golang/db/mock"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type EqCreateUserParamsMatcher struct {
	arg      db.CreateUserRequest
	password string
}

func (e EqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserRequest)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e EqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserRequest, password string) gomock.Matcher {
	return EqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(response *http.Response)
	}{
		{
			name: "CreateUserSuccess",
			body: types.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserRequest{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusCreated, response.StatusCode)
				requireBodyMatchUser(t, response.Body, user)
			},
		},
		{
			name: "InternalError",
			body: types.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "DuplicateUsername",
			body: types.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidUsername",
			body: types.CreateUserRequest{
				Username: "invalid-user#",
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: types.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "IncompleteReqBody",
			body: "akd",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "TooShortPassword",
			body: types.CreateUserRequest{
				Username: user.Username,
				Password: "123",
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
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

			server := NewServer(store)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			response, err := server.router.Test(request)
			require.NoError(t, err)

			tc.checkResponse(response)
		})
	}
}

// func TestGetUserAPI(t *testing.T) {
// 	user := randomUser()

// 	testCases := []struct {
// 		name          string
// 		body          interface{}
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, response *http.Response)
// 	}{
// 		{
// 			name: "GetUserSuccess",
// 			body: "akd",
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(user, nil)
// 			},
// 			checkResponse: func(t *testing.T, response *http.Response) {
// 				require.Equal(t, http.StatusOK, response.StatusCode)
// 				requireUserBodyMatch(t, response.Body, user)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: "akd",
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(t *testing.T, response *http.Response) {
// 				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
// 			},
// 		},
// 		{
// 			name: "UserNotFound",
// 			body: "akd",
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, response *http.Response) {
// 				require.Equal(t, http.StatusNotFound, response.StatusCode)
// 			},
// 		},
// 		{
// 			name: "InvalidParams",
// 			body: "a",
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, response *http.Response) {
// 				require.Equal(t, http.StatusBadRequest, response.StatusCode)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		ctrl := gomock.NewController(t)
// 		ctrl.Finish()

// 		store := mockdb.NewMockStore(ctrl)
// 		server := NewServer(store)
// 		tc.buildStubs(store)

// 		// get username parameter from body
// 		username := tc.body.(string)

// 		url := "/users/" + username

// 		request, err := http.NewRequest(http.MethodGet, url, nil)
// 		fmt.Printf("url: %v\n", url)
// 		require.NoError(t, err)

// 		response, err := server.router.Test(request)
// 		require.NoError(t, err)

//			tc.checkResponse(t, response)
//		}
//	}
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.GenerateRandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.GenerateRandomOwnerName(),
		HashedPassword: hashedPassword,
		FullName:       util.GenerateRandomOwnerName(),
		Email:          util.GenerateRandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body io.Reader, expected interface{}) {
	var userResponse struct {
		Data db.User `json:"data"`
	}
	err := json.NewDecoder(body).Decode(&userResponse)
	require.NoError(t, err)

	expectedUser := expected.(db.User)

	require.Equal(t, expectedUser.Username, userResponse.Data.Username)
	require.Equal(t, expectedUser.FullName, userResponse.Data.FullName)
	require.Equal(t, expectedUser.Email, userResponse.Data.Email)
	require.WithinDuration(t, expectedUser.CreatedAt, userResponse.Data.CreatedAt, time.Second)

}
