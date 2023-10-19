package gapi

import (
	"context"
	"database/sql"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := validateLoginUserRequest(req)
	if violations != nil {
		badRequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}
		statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

		statusDetails, err := statusInvalid.WithDetails(badRequest)
		if err != nil {
			return nil, statusInvalid.Err()
		}

		return nil, statusDetails.Err()
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found, try registering first")
		}
		return nil, status.Errorf(codes.Internal, "unexpected database error %v", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token")
	}

	mtd := server.extractMetadata(ctx)
	arg := db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ClientIp:     mtd.ClientIp,
		UserAgent:    mtd.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session")
	}

	resp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return resp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: err.Error(),
		})
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: err.Error(),
		})
	}

	return violations
}
