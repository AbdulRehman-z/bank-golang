package api

import (
	"testing"
	"time"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SYMMETRIC_KEY:         util.GenerateRandomString(32),
		ACCESS_TOKEN_DURATION: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
