package gapi

import (
	"context"
	"fmt"
	"main/database/db"
	"main/token"
	"main/util"
	"main/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := &util.ConfigDatabase{
		SecretKey:     util.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(store, config, taskDistributor)
	require.NoError(t, err)

	return server
}

func randomUser(t *testing.T) (*db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user := &db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return user, password
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
