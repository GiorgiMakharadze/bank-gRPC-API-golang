package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	"github.com/GiorgiMakharadze/bank-API-golang/token"
	"github.com/GiorgiMakharadze/bank-API-golang/util"
	"github.com/GiorgiMakharadze/bank-API-golang/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationHeaderBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
