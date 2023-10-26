package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.GenerateRandomString(8))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.GenerateRandomOwnerName(),
		HashedPassword: hashedPassword,
		FullName:       util.GenerateRandomOwnerName(),
		Email:          util.GenerateRandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserFullNameOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: util.GenerateRandomOwnerName(),
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEmpty(t, updatedUser)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt)
}

func TestUpdateUserPasswordOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)

	hashedPassword, err := util.HashPassword(util.GenerateRandomString(8))
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEmpty(t, updatedUser)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.NotEqual(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt)
}

func TestUpdateUserEmailOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: util.GenerateRandomEmail(),
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEmpty(t, updatedUser)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.NotEqual(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt)
}
