package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgareway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		if UserAgent := md.Get(userAgentHeader); len(UserAgent) > 0 {
			mtd.UserAgent = UserAgent[0]
		}

		if clientIps := md.Get(xForwardedForHeader); len(clientIps) > 0 {
			mtd.ClientIp = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIp = p.Addr.String()
	}

	return mtd
}
