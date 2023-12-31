package gapi

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	util2 "github.com/ngohoang211020/simplebank/gapi/util"
	pbuser "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *GrpcServer) CreateUser(ctx context.Context, req *pbuser.CreateUserRequest) (*pbuser.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pbuser.CreateUserResponse{
		User: util2.ConvertUser(user),
	}
	return rsp, nil
}

func (server *GrpcServer) UpdateUser(ctx context.Context, req *pbuser.UpdateUserRequest) (*pbuser.UpdateUserResponse, error) {
	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pbuser.UpdateUserResponse{
		User: util2.ConvertUser(user),
	}
	return rsp, nil
}
