package gapi

import (
	"context"
	"time"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/AbdulRehman-z/bank-golang/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
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

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password")
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			FullName:       req.GetFullname(),
			HashedPassword: hashedPassword,
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			payload := &worker.PayloadSendVerificationEmail{
				Username: user.Username,
			}

			err := server.taskDistributor.TaskSendVerificationEmail(ctx, payload, opts...)
			if err != nil {
				return status.Errorf(codes.Internal, "cannot send verification email")
			}

			return nil
		},
	}

	txUserResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" {
				return nil, status.Errorf(codes.InvalidArgument, "username already taken")
			}
		}
		return nil, status.Errorf(codes.Internal, "unexpected database error %v", err)
	}

	// opts := []asynq.Option{
	// 	asynq.MaxRetry(10),
	// 	asynq.ProcessIn(10 * time.Second),
	// 	asynq.Queue(worker.QueueCritical),
	// }

	// payload := &worker.PayloadSendVerificationEmail{
	// 	Username: txUserResult.Username,
	// }

	// err = server.taskDistributor.TaskSendVerificationEmail(ctx, payload, opts...)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "cannot send verification email")
	// }

	response := &pb.CreateUserResponse{
		User: &pb.User{
			Username:          txUserResult.User.Username,
			FullName:          txUserResult.User.FullName,
			Email:             txUserResult.User.Email,
			PasswordChangedAt: timestamppb.New(txUserResult.User.PasswordChangedAt),
			CreatedAt:         timestamppb.New(txUserResult.User.CreatedAt),
		},
	}

	return response, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: err.Error(),
		})
	}

	if err := ValidateFullName(req.GetFullname()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "fullname",
			Description: err.Error(),
		})
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: err.Error(),
		})
	}

	if err := ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email",
			Description: err.Error(),
		})
	}

	return violations
}
