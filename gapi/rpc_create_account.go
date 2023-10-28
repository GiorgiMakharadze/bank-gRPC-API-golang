package gapi

import (
	"context"

	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	"github.com/GiorgiMakharadze/bank-API-golang/pb"
	"github.com/GiorgiMakharadze/bank-API-golang/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	violations := ValidateCreateAccountRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	if payload.Username != req.Owner {
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission to create this account")
	}

	arg := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "account already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create account: %s", err)
	}

	rsp := &pb.CreateAccountResponse{
		Account: convertAccount(account),
	}
	return rsp, nil
}

func ValidateCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}
	return violations
}
