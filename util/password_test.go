package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	randomPassword := GenerateRandomString(8)
	hashedPassword, err := HashPassword(randomPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(randomPassword, hashedPassword)
	require.NoError(t, err)
	require.NotEqual(t, randomPassword, hashedPassword)

	wrongPassword := GenerateRandomString(8)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, "crypto/bcrypt: hashedPassword is not the hash of the given password")
}
