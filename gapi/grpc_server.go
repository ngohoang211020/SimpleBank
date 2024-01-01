package gapi

import (
	"fmt"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	pb "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/token"
	"github.com/ngohoang211020/simplebank/util"
)

// GrpcServer serves gRPC requests for our banking service.
type GrpcServer struct {
	config     *util.Configuration
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
}

// NewGrpcServer creates a new gRPC server
func NewGrpcServer(config *util.Configuration, store db.Store) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &GrpcServer{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	return server, nil
}
