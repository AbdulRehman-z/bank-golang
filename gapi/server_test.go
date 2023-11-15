package gapi

import (
	"testing"
	"time"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/AbdulRehman-z/bank-golang/worker"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {

	config := util.Config{
		SYMMETRIC_KEY:         util.GenerateRandomString(32),
		ACCESS_TOKEN_DURATION: time.Minute * 2,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
