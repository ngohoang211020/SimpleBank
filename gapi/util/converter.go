package util

import (
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	pb "github.com/ngohoang211020/simplebank/pb/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(userDb db.Users) *pb.User {
	return &pb.User{
		Username:          userDb.Username,
		FullName:          userDb.FullName,
		Email:             userDb.Email,
		PasswordChangedAt: timestamppb.New(userDb.PasswordChangedAt),
		CreatedAt:         timestamppb.New(userDb.CreatedAt),
	}
}
