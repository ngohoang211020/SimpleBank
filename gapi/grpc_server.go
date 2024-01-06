package gapi

import (
	"fmt"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	"github.com/ngohoang211020/simplebank/mail"
	pb "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/token"
	"github.com/ngohoang211020/simplebank/util"
	"github.com/ngohoang211020/simplebank/worker"
)

// GrpcServer serves gRPC requests for our banking service.
type GrpcServer struct {
	config     *util.Configuration
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
	distributor worker.TaskDistributor
	mail        mail.EmailSender
}

// NewGrpcServer creates a new gRPC server
func NewGrpcServer(config *util.Configuration, store db.Store, distributor worker.TaskDistributor, mail mail.EmailSender) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &GrpcServer{
		store:       store,
		tokenMaker:  tokenMaker,
		config:      config,
		distributor: distributor,
		mail:        mail,
	}
	return server, nil
}
