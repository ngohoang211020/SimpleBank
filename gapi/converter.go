package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	db "simplebank/db/sqlc"
	"simplebank/pb/user"
)

func convertUser(userDb db.Users) *user.User {
	return &user.User{
		Username:          userDb.Username,
		FullName:          userDb.FullName,
		Email:             userDb.Email,
		PasswordChangedAt: timestamppb.New(userDb.PasswordChangedAt),
		CreatedAt:         timestamppb.New(userDb.CreatedAt),
	}
}
