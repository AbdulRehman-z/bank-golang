package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/AbdulRehman-z/bank-golang/db/mock"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/AbdulRehman-z/bank-golang/worker"
	mockworker "github.com/AbdulRehman-z/bank-golang/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type EqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected EqCreateUserTxParamsMatcher) Matches(x any) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.AfterCreate(expected.user)

	return err == nil
}

func (e EqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return EqCreateUserTxParamsMatcher{
		arg:      arg,
		password: password,
		user:     user,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	// fmt.Println(user.Username)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
		checkResponse func(t *testing.T, response *pb.CreateUserResponse, err error)
	}{
		{
			name: "CreateUserSuccess",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				Fullname: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						FullName:       user.FullName,
						HashedPassword: user.HashedPassword,
						Username:       user.Username,
						Email:          user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateuserTxResult{User: user}, nil)

				payload := &worker.PayloadSendVerificationEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					TaskSendVerificationEmail(gomock.Any(), payload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, resp)
				// fmt.Printf("=====================================")
				// fmt.Printf("Expected: %v", user.Username)
				// fmt.Printf("Actual: %v", resp.User.Username)
				// fmt.Printf("=====================================")
				require.Equal(t, user.Username, resp.User.Username)
				require.Equal(t, user.Email, resp.User.Email)
				require.Equal(t, user.FullName, resp.User.FullName)
			},
		},
		// {
		// 	name: "InternalError",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		Fullname: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusInternalServerError, response.StatusCode)
		// 	},
		// },
		// {
		// 	name: "DuplicateUsername",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		Fullname: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, &pq.Error{Code: "23505"})
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		// 	},
		// },
		// {
		// 	name: "InvalidUsername",
		// 	body: &pb.CreateUserRequest{
		// 		Username: "invalid-user#",
		// 		Password: password,
		// 		Fullname: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		// 	},
		// },
		// {
		// 	name: "InvalidEmail",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		Fullname: user.FullName,
		// 		Email:    "invalid-email",
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		// 	},
		// },
		// {
		// 	name: "IncompleteReqBody",
		// 	body: &pb.CreateUserRequest{},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		// 	},
		// },
		// {
		// 	name: "TooShortPassword",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: "123",
		// 		Fullname: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(response *http.Response) {
		// 		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			taskCtrl := gomock.NewController(t)
			taskDistributor := mockworker.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store, taskDistributor)
			server := NewTestServer(t, store, taskDistributor)
			resp, err := server.CreateUser(context.Background(), tc.req)
			fmt.Println("=================================")
			fmt.Printf("Request: %v", tc.req.Username)
			fmt.Printf("Response: %v", resp.User.Username)
			fmt.Println("=================================")
			tc.checkResponse(t, resp, err)
		})
	}
}

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
