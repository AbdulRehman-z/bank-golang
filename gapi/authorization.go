package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/AbdulRehman-z/bank-golang/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authTypeBearer      = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization token is not provided")
	}

	accessToken := values[0]
	fields := strings.Fields(accessToken)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization token")
	}

	authType := strings.ToLower(fields[0])
	if authType != authTypeBearer {
		return nil, fmt.Errorf("authorization type is not supported")
	}

	payload, err := server.tokenMaker.VerifyToken(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid authorization token: %w", err)
	}

	return payload, nil
}
