package gapi

import (
	"context"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	pbuser "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GrpcServer) VerifyEmail(ctx context.Context, req *pbuser.VerifyEmailRequest) (*pbuser.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	arg := db.VerifyEmailTxParams{
		SecretCode: req.GetSecretCode(),
		EmailId:    req.GetEmailId(),
	}

	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pbuser.VerifyEmailResponse{
		IsVerified: result.User.IsVerifiedEmail,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *pbuser.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("emailId", err))
	}

	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secretCode", err))
	}

	return violations
}
