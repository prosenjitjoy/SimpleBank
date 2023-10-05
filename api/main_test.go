package api

import (
	"main/database/db"
	"main/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := &util.ConfigDatabase{
		SecretKey:     util.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}
