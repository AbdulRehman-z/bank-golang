package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, err
	}

	violations := validateUpdateUserRequest(req)
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

	if authPayload.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "permission denied: not allowed to update other user's")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullname(),
			Valid:  req.Fullname != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password")
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  req.Password != nil,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	respnse := &pb.UpdateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			Email:             user.Email,
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return respnse, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: err.Error(),
		})
	}

	if req.Fullname != nil {
		if err := ValidateFullName(req.GetFullname()); err != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "fullname",
				Description: err.Error(),
			})
		}
	}

	if req.Password != nil {
		if err := ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "password",
				Description: err.Error(),
			})
		}
	}

	if req.Email != nil {
		if err := ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "email",
				Description: err.Error(),
			})
		}
	}
	return violations
}
