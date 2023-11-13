package gapi

import (
	"context"
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {

	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		badReqest := &errdetails.BadRequest{
			FieldViolations: violations,
		}

		statusInvaid := status.New(codes.InvalidArgument, "invalid arguments")
		statusDetails, err := statusInvaid.WithDetails(badReqest)
		if err != nil {
			return nil, statusInvaid.Err()
		}

		return nil, statusDetails.Err()
	}

	verifyTx, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	response := &pb.VerifyEmailResponse{
		IsVerified: verifyTx.User.IsEmailVerified,
	}

	return response, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email_id",
			Description: err.Error(),
		})
	}

	if err := VlaidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "secret_code",
			Description: err.Error(),
		})
	}
	return violations
}
