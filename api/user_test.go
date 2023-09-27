package api

import (
	"bytes"
	"encoding/json"
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
			name: "OK",
			body: types.CreateUserParams{
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

func randomUser() db.User {
	return db.User{
		Username:       util.GenerateRandomOwnerName(),
		FullName:       util.GenerateRandomOwnerName(),
		Email:          util.GenerateRandomEmail(),
		HashedPassword: "secret",
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
