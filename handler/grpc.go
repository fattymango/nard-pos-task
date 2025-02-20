package handler

import (
	"context"
	"fmt"
	"multitenant/internal/engine"
	"multitenant/model"
	"multitenant/pkg/config"
	"multitenant/pkg/logger"
	pb "multitenant/proto/multitenant"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MultiTenantRPCServer struct {
	pb.UnimplementedMultiTenantServer
	config *config.Config
	logger *logger.Logger
	engine *engine.Engine
	server *grpc.Server
	ctx    context.Context
}

func NewMultiTenantRPCServer(ctx context.Context, config *config.Config, logger *logger.Logger, engine *engine.Engine) *MultiTenantRPCServer {
	return &MultiTenantRPCServer{
		config: config,
		logger: logger,
		engine: engine,
		ctx:    ctx,
	}
}

func (s *MultiTenantRPCServer) CreateTransaction(ctx context.Context, req *pb.CrtTransaction) (*pb.TransactionResponse, error) {
	// s.logger.Debugf("Received CreateTransaction request: %v", req)

	tx := &model.Transaction{
		TenantID:     req.TenantId,
		BranchID:     req.BranchId,
		ProductID:    req.ProductId,
		QuantitySold: req.QuantitySold,
		PricePerUnit: req.PricePerUnit,
		Status:       model.TransactionStatusPending,
	}
	// Create Transaction
	err := s.engine.CreateTransaction(tx)
	if err != nil {
		return &pb.TransactionResponse{
			Message: fmt.Sprintf("Failed to create transaction: %v", err),
			Success: false,
		}, nil
	}

	return &pb.TransactionResponse{
		Message: "Transaction Created",
		Success: true,
	}, nil

}

func (s *MultiTenantRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	creds := insecure.NewCredentials()
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(s.rpcDurationInterceptor))

	s.server = grpcServer

	pb.RegisterMultiTenantServer(grpcServer, s)

	s.logger.Infof("gRPC server listening on %s", s.config.GRPC.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *MultiTenantRPCServer) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.server.GracefulStop()
}
