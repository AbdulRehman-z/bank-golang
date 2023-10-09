package token

import (
	"testing"
	"time"

	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoCreateToken(t *testing.T) {
	symmetricKey := util.GenerateRandomString(32)

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := util.GenerateRandomOwnerName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.Id)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoExpiredToken(t *testing.T) {
	symmetricKey := util.GenerateRandomString(32)

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := util.GenerateRandomOwnerName()
	duration := -time.Minute

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
