package client

import (
	"context"
	"google.golang.org/grpc"
	"simplebank/common/constants"
	"simplebank/pb/user"
	"time"
)

type AuthClient struct {
	authServiceClient user.AuthServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	authClient := user.NewAuthServiceClient(cc)
	return &AuthClient{
		authServiceClient: authClient,
	}
}

func (service *AuthClient) Login(req *user.LoginUserRequest) (res *user.LoginUserResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.GrpcTimeoutInSecs*time.Second)
	defer cancel()
	res, err = service.authServiceClient.Login(ctx, req)
	return
}
