package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcUserAgentHeader        = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetaData(ctx context.Context) *Metadata {
	metaData := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userAgents := md.Get(grpcGatewayUserAgentHeader)
		if len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		userAgents = md.Get(grpcUserAgentHeader)
		if len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		clientIPs := md.Get(xForwardedForHeader)
		if len(clientIPs) > 0 {
			metaData.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		metaData.ClientIP = p.Addr.String()
	}

	return metaData
}
