package handler

import (
	"context"

	"google.golang.org/grpc"
)

// gRPC Interceptor to log the duration of the RPC call
func (s *MultiTenantRPCServer) rpcDurationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// start := time.Now()
	resp, err := handler(ctx, req)
	// duration := time.Since(start)

	// Log the duration of the RPC call
	// s.logger.Infof("RPC %s took %v", info.FullMethod, duration)

	return resp, err
}
