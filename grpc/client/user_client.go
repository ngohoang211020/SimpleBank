package client

import (
	"context"
	"google.golang.org/grpc"
	"simplebank/common/constants"
	"simplebank/pb/user"
	"time"
)

type UserClient struct {
	userServiceClient user.UserServiceClient
}

func NewUserClient(cc *grpc.ClientConn) *UserClient {
	userServiceClient := user.NewUserServiceClient(cc)
	return &UserClient{
		userServiceClient: userServiceClient,
	}
}

func (service *UserClient) CreateUser(req *user.CreateUserRequest) (res *user.CreateUserResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.GrpcTimeoutInSecs*time.Second)
	defer cancel()
	res, err = service.userServiceClient.CreateUser(ctx, req)
	return
}

func (service *UserClient) UpdateUser(req *user.UpdateUserRequest) (res *user.UpdateUserResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.GrpcTimeoutInSecs*time.Second)
	defer cancel()
	res, err = service.userServiceClient.UpdateUser(ctx, req)
	return
}

func (service *UserClient) VerifyEmail(req *user.VerifyEmailRequest) (res *user.VerifyEmailResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.GrpcTimeoutInSecs*time.Second)
	defer cancel()
	res, err = service.userServiceClient.VerifyEmail(ctx, req)
	return
}
