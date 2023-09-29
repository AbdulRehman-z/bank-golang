package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoCreateToken(t *testing.T) {
	symmetricKey := GenerateRandomString(32)

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := GenerateRandomOwnerName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.Id)
	require.WithinDuration(t, issuedAt, time.Unix(payload.IssuedAt, 0), time.Second)
	require.WithinDuration(t, expiredAt, time.Unix(payload.ExpiredAt, 0), time.Second)
}

func TestPasetoExpiredToken(t *testing.T) {
	symmetricKey := GenerateRandomString(32)

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := GenerateRandomOwnerName()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
